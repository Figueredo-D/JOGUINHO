package audio

import (
	"bytes"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

// TrackID identifica cada faixa de trilha sonora do jogo.
type TrackID int

const (
	TrackNone TrackID = iota
	TrackMenu
	TrackExploration
	TrackBoss
)

var (
	bgmBank      map[TrackID][]byte
	currentTrack TrackID
	bgmPlayer    *audio.Player
)

func initBGM() {
	bgmBank = make(map[TrackID][]byte)
	bgmBank[TrackMenu] = genMenuBGM()
	bgmBank[TrackExploration] = genExplorationBGM()
	bgmBank[TrackBoss] = genBossBGM()
}

// PlayBGM inicia ou alterna para a faixa musical solicitada em loop infinito.
func PlayBGM(id TrackID) {
	if ctx == nil || bgmBank == nil {
		return
	}
	if id == currentTrack && bgmPlayer != nil && bgmPlayer.IsPlaying() {
		return
	}

	// Pausa player anterior se existir
	if bgmPlayer != nil {
		_ = bgmPlayer.Close()
		bgmPlayer = nil
	}
	currentTrack = id

	if id == TrackNone {
		return
	}

	data, ok := bgmBank[id]
	if !ok || len(data) == 0 {
		return
	}

	loop := audio.NewInfiniteLoop(bytes.NewReader(data), int64(len(data)))
	p, err := ctx.NewPlayer(loop)
	if err != nil {
		return
	}
	bgmPlayer = p
	updateBGMVolume()
	bgmPlayer.Play()
}

func updateBGMVolume() {
	if bgmPlayer == nil {
		return
	}
	mu.Lock()
	vol := masterVolume * musicVolume
	mu.Unlock()
	bgmPlayer.SetVolume(vol * 0.25) // BGM suave para não encobrir os SFX
}

// Onda triangular suave (ideal para baixos e pads medievais/chiptune)
func triangleWave(freq, t float64) float64 {
	phase := math.Mod(freq*t, 1.0)
	if phase < 0.25 {
		return phase * 4.0
	} else if phase < 0.75 {
		return 2.0 - phase*4.0
	}
	return phase*4.0 - 4.0
}

// Onda pulso 25% (som clássico 8-bit chiptune)
func pulseWave(freq, t, duty float64) float64 {
	if math.Mod(freq*t, 1.0) < duty {
		return 1.0
	}
	return -1.0
}

// 1. Música do Menu Principal: Tom misterioso, calmo e solene (Ré menor)
func genMenuBGM() []byte {
	// Progressão solene de 8 segundos: Dm - Bb - Gm - A
	// Notas em Hz: D3=146.83, F3=174.61, A3=220.00, Bb3=233.08, G3=196.00, E3=164.81, C#4=277.18, D4=293.66
	chords := [][]float64{
		{146.83, 220.00, 293.66, 349.23}, // Dm
		{116.54, 174.61, 233.08, 293.66}, // Bb
		{98.00, 196.00, 233.08, 293.66},  // Gm
		{110.00, 164.81, 220.00, 277.18}, // A
	}

	duration := 8.0
	return synthBuffer(duration, func(t, p float64) float64 {
		chordIdx := int(t/2.0) % len(chords)
		chord := chords[chordIdx]
		localT := math.Mod(t, 2.0)

		// Pad grave calmo e solene
		pad := math.Sin(2*math.Pi*chord[0]*t)*0.35 + triangleWave(chord[1], t)*0.25

		// Arpejo lento misterioso (uma nota a cada 0.5s)
		noteIdx := int(localT/0.5) % len(chord)
		noteT := math.Mod(localT, 0.5)
		env := math.Exp(-noteT * 3.5)
		arp := triangleWave(chord[noteIdx]*2.0, t) * env * 0.35

		return (pad + arp) * 0.75
	})
}

// 2. Faixa de Exploração / Hordas: Baixo pulsante constante + arpejos sombrios de masmorra
func genExplorationBGM() []byte {
	// 120 BPM -> 0.25s por colcheia (16 passos por 4s, loop de 8s)
	bassNotes := []float64{
		73.42, 73.42, 87.31, 73.42, 98.00, 87.31, 82.41, 73.42, // D2, F2, G2, E2
		73.42, 73.42, 110.00, 98.00, 87.31, 82.41, 69.30, 82.41, // A2, C#2
	}
	leadNotes := []float64{
		293.66, 349.23, 440.00, 349.23, 392.00, 349.23, 329.63, 293.66,
		293.66, 440.00, 523.25, 440.00, 392.00, 349.23, 329.63, 277.18,
	}

	stepDur := 0.25
	duration := float64(len(bassNotes)) * stepDur * 2.0 // 8.0 segundos

	return synthBuffer(duration, func(t, p float64) float64 {
		step := int(t/stepDur) % len(bassNotes)
		stepT := math.Mod(t, stepDur)

		// Baixo pulsante rítmico
		bEnv := math.Exp(-stepT * 7.0)
		bass := pulseWave(bassNotes[step], t, 0.25) * bEnv * 0.45

		// Arpejo chiptune de masmorra
		lEnv := math.Exp(-stepT * 5.0)
		lead := triangleWave(leadNotes[step], t) * lEnv * 0.35

		// Sub-pulso de passo (bumbo leve nos tempos pares)
		kick := 0.0
		if step%2 == 0 && stepT < 0.09 {
			kEnv := math.Exp(-stepT * 35.0)
			kick = math.Sin(2*math.Pi*(110.0-stepT*600.0)*stepT) * kEnv * 0.40
		}

		return (bass + lead + kick) * 0.70
	})
}

// 3. Faixa de Boss Fight: Ritmo acelerado, percussão agressiva e melodia urgente (150 BPM)
func genBossBGM() []byte {
	rng := rand.New(rand.NewSource(777))
	// 150 BPM -> 0.1s por semicolcheia (16 passos = 1.6s; loop de 6.4s)
	bassSeq := []float64{
		73.42, 73.42, 146.83, 73.42, 87.31, 73.42, 146.83, 98.00,
		73.42, 73.42, 146.83, 73.42, 110.00, 98.00, 87.31, 82.41,
	}
	leadSeq := []float64{
		440.00, 523.25, 587.33, 440.00, 659.25, 587.33, 523.25, 440.00,
		587.33, 698.46, 880.00, 698.46, 659.25, 587.33, 554.37, 440.00,
	}

	stepDur := 0.125
	duration := float64(len(bassSeq)) * stepDur * 2.0 // 4.0 segundos por ciclo

	// Pré-gera um pequeno ruído determinístico para a caixa (snare)
	noiseSamples := make([]float64, int(float64(SampleRate)*stepDur)+1)
	for i := range noiseSamples {
		noiseSamples[i] = rng.Float64()*2.0 - 1.0
	}

	return synthBuffer(duration, func(t, p float64) float64 {
		step := int(t/stepDur) % len(bassSeq)
		stepT := math.Mod(t, stepDur)

		// Baixo agressivo contínuo
		bEnv := math.Exp(-stepT * 9.0)
		bass := pulseWave(bassSeq[step], t, 0.5) * bEnv * 0.45

		// Melodia urgente aguda
		lEnv := math.Exp(-stepT * 6.0)
		lead := pulseWave(leadSeq[step], t, 0.25) * lEnv * 0.30

		// Percussão: Bumbo nos passos 0,4,8,12 e Caixa nos passos 2,6,10,14
		perc := 0.0
		if step%4 == 0 && stepT < 0.09 {
			kEnv := math.Exp(-stepT * 40.0)
			perc = math.Sin(2*math.Pi*(150.0-stepT*900.0)*stepT) * kEnv * 0.60
		} else if step%4 == 2 && stepT < 0.07 {
			sIdx := int(stepT * SampleRate)
			if sIdx < len(noiseSamples) {
				sEnv := math.Exp(-stepT * 45.0)
				perc = noiseSamples[sIdx] * sEnv * 0.45
			}
		}

		return (bass + lead + perc) * 0.72
	})
}

package audio

import (
	"math"
	"math/rand"
	"sync"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const SampleRate = 44100

// SoundID identifica cada efeito sonoro sintetizado.
type SoundID int

const (
	SFXSlash SoundID = iota
	SFXArrow
	SFXMagic
	SFXShotgun
	SFXRevolver
	SFXHit
	SFXExplosion
	SFXBossAlert
	SFXPickup
	SFXDash
	SFXLightning
	SFXRelic
	SFXShieldReflect
)

var (
	ctx          *audio.Context
	initOnce     sync.Once
	soundBank    map[SoundID][]byte
	masterVolume = 0.80
	sfxVolume    = 0.80
	musicVolume  = 0.60
	lastPlayMap  = make(map[SoundID]int)
	globalFrame  int
	mu           sync.Mutex
)

// Init inicializa o contexto de áudio e pré-sintetiza todos os efeitos e faixas musicais em memória.
func Init() {
	initOnce.Do(func() {
		ctx = audio.NewContext(SampleRate)
		soundBank = make(map[SoundID][]byte)

		soundBank[SFXSlash] = genSlash()
		soundBank[SFXArrow] = genArrow()
		soundBank[SFXMagic] = genMagic()
		soundBank[SFXShotgun] = genShotgun()
		soundBank[SFXRevolver] = genRevolver()
		soundBank[SFXHit] = genHit()
		soundBank[SFXExplosion] = genExplosion()
		soundBank[SFXBossAlert] = genBossAlert()
		soundBank[SFXPickup] = genPickup()
		soundBank[SFXDash] = genDash()
		soundBank[SFXLightning] = genLightning()
		soundBank[SFXRelic] = genRelic()
		soundBank[SFXShieldReflect] = genReflect()

		initBGM()
	})
}

// Tick avança o contador de frames para limitar sobreposição excessiva de sons idênticos.
func Tick() {
	mu.Lock()
	globalFrame++
	mu.Unlock()
}

func clampPercent(p int) float64 {
	if p < 0 {
		return 0
	}
	if p > 100 {
		return 1
	}
	return float64(p) / 100.0
}

// SetMasterVolume ajusta o volume mestre (0 a 100) e atualiza a música ativa.
func SetMasterVolume(percent int) {
	mu.Lock()
	masterVolume = clampPercent(percent)
	mu.Unlock()
	updateBGMVolume()
}

// SetSFXVolume ajusta o volume dos efeitos sonoros (0 a 100).
func SetSFXVolume(percent int) {
	mu.Lock()
	sfxVolume = clampPercent(percent)
	mu.Unlock()
}

// SetMusicVolume ajusta o volume da trilha sonora (0 a 100) e atualiza a faixa em reprodução.
func SetMusicVolume(percent int) {
	mu.Lock()
	musicVolume = clampPercent(percent)
	mu.Unlock()
	updateBGMVolume()
}

// SetVolume ajusta o volume mestre (retrocompatibilidade).
func SetVolume(percent int) {
	SetMasterVolume(percent)
}

// Play toca um efeito sonoro pré-sintetizado sem alocar novos buffers PCM.
func Play(id SoundID) {
	if ctx == nil || soundBank == nil {
		return
	}
	mu.Lock()
	effVol := masterVolume * sfxVolume
	lastFrame := lastPlayMap[id]
	if globalFrame-lastFrame < 3 && id != SFXExplosion && id != SFXPickup {
		mu.Unlock()
		return
	}
	lastPlayMap[id] = globalFrame
	mu.Unlock()

	if effVol <= 0.01 {
		return
	}
	data, ok := soundBank[id]
	if !ok || len(data) == 0 {
		return
	}

	p := ctx.NewPlayerFromBytes(data)
	p.SetVolume(effVol * 0.45)
	p.Play()
}

// synthBuffer gera amostras PCM 16-bit estéreo (4 bytes/sample).
func synthBuffer(durationSec float64, fn func(t, progress float64) float64) []byte {
	numSamples := int(float64(SampleRate) * durationSec)
	buf := make([]byte, numSamples*4)
	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(SampleRate)
		progress := float64(i) / float64(numSamples)
		val := fn(t, progress)
		if val > 1.0 {
			val = 1.0
		} else if val < -1.0 {
			val = -1.0
		}
		sample := int16(val * 32767.0)
		idx := i * 4
		buf[idx] = byte(sample)
		buf[idx+1] = byte(sample >> 8)
		buf[idx+2] = byte(sample)
		buf[idx+3] = byte(sample >> 8)
	}
	return buf
}

// Corte seco de espada (ruído filtrado + queda rápida de frequência)
func genSlash() []byte {
	rng := rand.New(rand.NewSource(101))
	return synthBuffer(0.14, func(t, p float64) float64 {
		env := math.Pow(1.0-p, 1.8)
		freq := 480.0 - p*340.0
		tone := math.Sin(2 * math.Pi * freq * t)
		noise := rng.Float64()*2.0 - 1.0
		return (tone*0.35 + noise*0.65) * env
	})
}

// Disparo rápido de arco e flecha
func genArrow() []byte {
	rng := rand.New(rand.NewSource(102))
	return synthBuffer(0.11, func(t, p float64) float64 {
		env := math.Pow(1.0-p, 2.2)
		freq := 680.0 - p*450.0
		twang := math.Sin(2 * math.Pi * freq * t)
		air := (rng.Float64()*2.0 - 1.0) * 0.45
		return (twang*0.6 + air) * env
	})
}

// Bipe agudo / onda mágica para feitiços
func genMagic() []byte {
	return synthBuffer(0.16, func(t, p float64) float64 {
		env := math.Sin(p * math.Pi)
		freq := 740.0 + math.Sin(p*20.0)*120.0 + p*260.0
		wave := math.Sin(2 * math.Pi * freq * t)
		harm := math.Sin(2*math.Pi*(freq*1.5)*t) * 0.3
		return (wave + harm) * env * 0.7
	})
}

// Estampido grave de espingarda calibre 12
func genShotgun() []byte {
	rng := rand.New(rand.NewSource(103))
	return synthBuffer(0.24, func(t, p float64) float64 {
		env := math.Pow(1.0-p, 1.6)
		bass := math.Sin(2*math.Pi*(140.0-p*100.0)*t) * 0.65
		blast := (rng.Float64()*2.0 - 1.0) * 0.85
		return (bass + blast) * env
	})
}

// Tiro seco de revólver do Xerife
func genRevolver() []byte {
	rng := rand.New(rand.NewSource(104))
	return synthBuffer(0.13, func(t, p float64) float64 {
		env := math.Pow(1.0-p, 2.5)
		snap := math.Sin(2*math.Pi*(320.0-p*220.0)*t) * 0.6
		crack := (rng.Float64()*2.0 - 1.0) * 0.75
		return (snap + crack) * env
	})
}

// Impacto curto ao acertar inimigos
func genHit() []byte {
	rng := rand.New(rand.NewSource(105))
	return synthBuffer(0.07, func(t, p float64) float64 {
		env := math.Pow(1.0-p, 2.0)
		thud := math.Sin(2*math.Pi*(220.0-p*140.0)*t) * 0.7
		n := (rng.Float64()*2.0 - 1.0) * 0.35
		return (thud + n) * env
	})
}

// Ruído grave de explosão / morte de monstro / habilidades em área
func genExplosion() []byte {
	rng := rand.New(rand.NewSource(106))
	return synthBuffer(0.32, func(t, p float64) float64 {
		env := math.Pow(1.0-p, 1.4)
		sub := math.Sin(2*math.Pi*(95.0-p*65.0)*t) * 0.75
		rumble := (rng.Float64()*2.0 - 1.0) * 0.65
		return (sub + rumble) * env
	})
}

// Sirene / Alerta estridente quando o Boss desperta
func genBossAlert() []byte {
	return synthBuffer(0.65, func(t, p float64) float64 {
		env := 1.0
		if p < 0.1 {
			env = p / 0.1
		} else if p > 0.8 {
			env = (1.0 - p) / 0.2
		}
		// Onda quadrada / dente de serra dramática oscilando
		freq := 320.0 + math.Sin(p*math.Pi*4.0)*140.0
		sq := 1.0
		if math.Sin(2*math.Pi*freq*t) < 0 {
			sq = -1.0
		}
		brass := math.Sin(2*math.Pi*(freq*0.5)*t) * 0.5
		return (sq*0.4 + brass) * env * 0.7
	})
}

// Glin/blip metálico sutil ao coletar almas ou itens
func genPickup() []byte {
	return synthBuffer(0.14, func(t, p float64) float64 {
		env := math.Pow(1.0-p, 1.5)
		freq := 880.0
		if p > 0.45 {
			freq = 1320.0 // Salto de quinta justa cristalino
		}
		return math.Sin(2*math.Pi*freq*t) * env * 0.65
	})
}

// Sopro rápido de esquiva / dash
func genDash() []byte {
	rng := rand.New(rand.NewSource(107))
	return synthBuffer(0.15, func(t, p float64) float64 {
		env := math.Sin(p * math.Pi)
		whoosh := (rng.Float64()*2.0 - 1.0) * 0.5
		sweep := math.Sin(2*math.Pi*(260.0+p*220.0)*t) * 0.25
		return (whoosh + sweep) * env
	})
}

// Chiado elétrico de Faísca Eletrostática
func genLightning() []byte {
	rng := rand.New(rand.NewSource(108))
	return synthBuffer(0.14, func(t, p float64) float64 {
		env := math.Pow(1.0-p, 1.8)
		buzz := 1.0
		if math.Sin(2*math.Pi*(950.0-p*400.0)*t) < 0 {
			buzz = -1.0
		}
		crackle := rng.Float64()*2.0 - 1.0
		return (buzz*0.4 + crackle*0.5) * env * 0.6
	})
}

// Acorde sagrado ao escolher uma relíquia
func genRelic() []byte {
	return synthBuffer(0.35, func(t, p float64) float64 {
		env := math.Pow(1.0-p, 1.3)
		c1 := math.Sin(2 * math.Pi * 523.25 * t)
		c2 := math.Sin(2 * math.Pi * 659.25 * t)
		c3 := math.Sin(2 * math.Pi * 783.99 * t)
		return (c1 + c2 + c3) / 3.0 * env * 0.8
	})
}

// Som metálico de reflexão no escudo
func genReflect() []byte {
	return synthBuffer(0.16, func(t, p float64) float64 {
		env := math.Pow(1.0-p, 1.7)
		freq := 1100.0 + p*400.0
		return math.Sin(2*math.Pi*freq*t) * env * 0.75
	})
}

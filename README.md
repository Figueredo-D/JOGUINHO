# As Catacumbas de Dorfieegu — Versão 1.0 Oficial

Um jogo roguelite de ação 2D desenvolvido nativamente em **Go** com a engine gráfica **Ebitengine v2**. Ambientado em um universo de **Fantasia Medieval e Velho Oeste Sombrio**, conta com campanha completa de **15 ondas**, trilha sonora chiptune dinâmica, áudio sintetizado, sistema de relíquias, eventos climáticos e gráficos em pixel art 100% procedurais (sem dependência de arquivos externos de imagem ou som).

---

## Estrutura do Projeto

```
├── cmd/
│   └── game/
│       └── main.go               # Janela ("As Catacumbas de Dorfieegu") e ícone procedural
├── internal/
│   ├── audio/
│   │   ├── synth.go              # Mixer (Master/SFX/Music) e SFX PCM procedurais
│   │   └── bgm.go                # Trilha sonora chiptune dinâmica (Menu, Exploração e Boss)
│   ├── world/
│   │   ├── arena.go              # Gerador procedural de arenas e biomas texturizados
│   │   └── weather.go            # Eventos climáticos (Névoa, Chuva Ácida e Lua Carmesim)
│   ├── progression/
│   │   ├── save.go               # Persistência (save.json): Almas, Upgrades, Conquistas e Opções
│   │   └── save_test.go          # Testes unitários de save/load e progressão
│   ├── sprites/
│   │   ├── sprites.go            # Pixel Art base (Heróis, Inimigos e Itens)
│   │   └── sprites_extra.go      # Pixel Art do Atirador Pesado, Xerife, Morcego e Chefe
│   ├── collision/
│   │   ├── aabb.go               # Detecção de colisões 2D (AABB)
│   │   └── aabb_test.go          # Testes unitários de colisão
│   ├── ui/
│   │   ├── button.go             # Botões estilizados medievais
│   │   └── bar.go                # Barras dinâmicas de progresso
│   ├── entities/
│   │   ├── player.go             # 5 Classes, Dash com after-image, I-frames e Relíquias
│   │   ├── relic.go              # Catálogo das 6 Relíquias Passivas acumuláveis
│   │   ├── boss.go               # Chefes (Onda 5, 10) e O Soberano de Dorfieegu (Onda 15)
│   │   ├── enemy.go              # 5 Arquétipos animados + flash branco de dano + knockback
│   │   ├── skills.go             # Habilidades ativas (Escudo, Chuva de Flechas, Pilar de Fogo)
│   │   ├── vfx.go                # Partículas, bolhas, ondas de choque e rachaduras
│   │   ├── drop.go               # Drops benignos/malignos e Fragmentos de Alma
│   │   ├── projectile.go         # Projéteis balísticos, teleguiados e perfurantes
│   │   └── slash.go              # Arco cortante Melee
│   ├── scenes/
│   │   ├── scene.go              # Interface Scene e SceneNavigator
│   │   ├── menu.go               # Menu Oficial (Logotipo Versão 1.0 e brasas ascendentes)
│   │   ├── how_to_play.go        # Guia de Comandos ("Como Jogar")
│   │   ├── sanctuary.go          # Santuário das Almas (Upgrades Permanentes das 5 classes)
│   │   ├── achievements.go       # Salão das Conquistas
│   │   ├── options.go            # Sliders de Volume (Master/SFX/BGM) e Tela Cheia
│   │   ├── class_select.go       # Seleção das 5 classes
│   │   ├── gameplay.go           # Loop completo (Screen Shake, Hitstop, Pausa com Ficha da Run)
│   │   ├── gameplay_test.go      # Testes unitários e de regressão
│   │   ├── victory.go            # Tela de Vitória Triunfal com relatório e bônus 1.5x
│   │   └── game_over.go          # Cripta da derrota com resgate de almas
│   └── game/
│       └── game.go               # Implementação de ebiten.Game e carregamento de opções
├── go.mod
└── go.sum
```

---

## Novidades da Versão 1.0 Oficial

- **Campanha Fechada em 15 Ondas & Vitória:**
  - *Onda 5:* **O Carrasco das Catacumbas** (Mini-chefe)
  - *Onda 10:* **O Golem de Almas** (Chefe Intermediário)
  - *Onda 15:* **O Soberano de Dorfieegu** (Chefe Final em 3 Fases com Feixes de Lasers Etéreos, Chuvas de Projéteis Radiais e Investidas Encadeadas).
  - Ao derrotar o Soberano na Onda 15, todos os projéteis inimigos viram almas, aplica-se um **Bônus de Vitória 1.5x** nas essências resgatadas e abre-se a tela de **Vitória: Catacumbas Purificadas** com relatório completo.
- **Game Feel & "Juice":** Screen Shake exponencial na câmera de mundo (`worldCanvas`), Hitstop (congelamento de 2 a 5 frames em impactos pesados), Flash Branco puro por frame de dano em todas as entidades e Knockback físico direcional.
- **Trilha Sonora Dinâmica & Mixer Persistido:** 3 músicas sintetizadas em loop contínuo (*Menu Solene*, *Exploração de Masmorra*, *Boss Fight Agressiva*) e sliders individuais para **Master**, **SFX** e **BGM**, salvos automaticamente em `save.json`.
- **Menu de Pausa Completo (`ESC` / `P`):** Exibe Ficha da Run atual com todas as Relíquias Passivas ativas e acesso inline às Opções de Áudio/Vídeo, Reiniciar Partida ou Voltar ao Menu Principal.

---

## Controles

- **W, A, S, D** ou **Setas**: Movimentação.
- **Mouse**: Mira e Clique Esquerdo para Ataque Básico.
- **Espaço** ou **Shift**: **Esquiva / Dash** (com I-Frames).
- **Tecla E** ou **Clique Direito**: Habilidade Especial Ativa da Classe.
- **Tecla R**: Recarregar Tambor (exclusivo do Xerife).
- **ESC** ou **P**: Pausar / Ficha da Run.

---

## Como Executar

```powershell
go run ./cmd/game
```

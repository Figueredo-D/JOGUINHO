# Dungeon Survivors (Roguelite 2D em Go)

Um jogo de sobrevivência roguelite 2D desenvolvido nativamente em **Go** com a engine gráfica **Ebitengine v2**, ambientado em um universo de **Fantasia Medieval** com gráficos em pixel art procedural em código (sem dependência de assets externos).

---

## Estrutura do Projeto

```
├── cmd/
│   └── game/
│       └── main.go           # Ponto de entrada da aplicação desktop
├── internal/
│   ├── progression/
│   │   ├── save.go           # Sistema de persistência meta (save.json) e custos de upgrade
│   │   └── save_test.go      # Testes unitários de save/load e progressão
│   ├── sprites/
│   │   └── sprites.go        # Gerador procedural de Pixel Art (Heróis, Inimigos, Drops e Almas)
│   ├── collision/
│   │   ├── aabb.go           # Detecção de colisões 2D (AABB)
│   │   └── aabb_test.go      # Testes unitários de colisão
│   ├── ui/
│   │   ├── button.go         # Botões estilizados medievais (carvalho, ferro e rebites)
│   │   └── bar.go            # Barras dinâmicas de progresso (HP, Mana, Cooldown)
│   ├── entities/
│   │   ├── player.go         # Herói com stats meta, postura de escudo e atributos por classe
│   │   ├── skills.go         # Efeitos de Habilidades Especiais (Chuva de Flechas, Pilar de Fogo)
│   │   ├── enemy.go          # 4 Arquétipos de monstros (Comum, Corredor, Tanque e Atirador)
│   │   ├── drop.go           # Drops de buffs/debuffs e Fragmentos de Alma meta
│   │   ├── projectile.go     # Projéteis do jogador e atiradores (com suporte a reflexão)
│   │   └── slash.go          # Arco cortante do ataque Melee de espada
│   ├── scenes/
│   │   ├── scene.go          # Interface Scene e SceneNavigator
│   │   ├── menu.go           # Menu Principal Medieval (tochas, brasão e atalho ao Santuário)
│   │   ├── sanctuary.go      # Santuário das Almas (Árvore de Upgrades Meta por classe)
│   │   ├── options.go        # Configurações do Reino (Fullscreen on/off, volume)
│   │   ├── class_select.go   # Seleção das 3 classes com prévia em Pixel Art
│   │   ├── gameplay.go       # Loop do combate, habilidades especiais ativas e pausa com ESC
│   │   └── game_over.go      # Cripta da derrota com resgate de almas e atalho de evolução
│   └── game/
│       └── game.go           # Implementação de ebiten.Game (Update, Draw, Layout)
├── go.mod
└── go.sum
```

---

## Sistema Meta-Roguelite (Santuário das Almas)

- **Moeda Meta (Fragmentos de Alma):** Monstros eliminados têm chance de dropar fragmentos em formato de cristal púrpura etéreo durante a batalha. Ao final da partida, a pontuação também é convertida em fragmentos extras.
- **Persistência Automática:** Salvo localmente em `save.json` e recarregado a cada inicialização.
- **Árvore de Upgrades Permanentes por Classe:**
  - **Guerreiro:**
    1. *Vitalidade Titânica:* +25 HP Máximo por nível.
    2. *Lâmina Dilacerante:* +12% Dano Melee e +6px de alcance do corte por nível.
    3. *Bastião Divino (Escudo):* Reduz o tempo de recarga da Postura de Escudo e amplia o dano refletido.
  - **Arqueiro:**
    1. *Passo do Vento:* +0.35 de Velocidade de Movimento permanente por nível.
    2. *Disparo Voraz:* +10% de Cadência de ataque e +1.2 na velocidade das flechas por nível.
    3. *Chuva Devastadora:* +15% de Área e +20% de Dano da Chuva de Flechas por nível.
  - **Mago:**
    1. *Fonte Mística:* +20 de Mana Máxima e +3.5 MP/s de regeneração por nível.
    2. *Canalização Arcana:* +15% de Dano Mágico em todos os feitiços por nível.
    3. *Cataclismo de Chamas:* +15% de Raio e +20 frames de duração do Pilar de Fogo por nível.

---

## Habilidades Especiais Ativas (`E` ou `Botão Direito do Mouse`)

Cada classe conta com uma habilidade própria e barra de recarga no topo do HUD:

- **Guerreiro — Postura de Escudo:**
  - Levanta uma barreira frontal luminosa e inquebrável por ~1.5 segundos.
  - Anula todo dano recebido frontalmente.
  - **Reflete projéteis inimigos** (como os tiros dos atiradores) de volta na direção oposta, multiplicando seu dano em mais de 2x e atingindo os próprios monstros!
- **Arqueiro — Chuva de Flechas:**
  - Dispara uma salva aos céus que cai velozmente em uma área circular selecionada pelo cursor.
  - Causa dano contínuo a todos os inimigos no raio e aplica **desaceleração (slow)** por 2 segundos.
- **Mago — Pilar de Fogo:**
  - Invoca uma erupção colossal de magma e labaredas no ponto do cursor.
  - Causa dano contínuo intenso em área, consumindo mana dedicado.

---

## Controles

- **W, A, S, D** ou **Setas**: Movimentação.
- **Mouse (Cursor)**: Mira.
- **Clique Esquerdo** ou **Espaço**: Ataque Básico da classe.
- **Tecla E** ou **Clique Direito**: **Habilidade Especial Ativa**.
- **ESC**: Menu de Pausa Medieval (Retomar, Tela Cheia ou Voltar ao Menu).

---

## Como Executar

```powershell
go run ./cmd/game
```

# Dungeon Survivors (Ebitengine 2D em Go)

Um jogo de sobrevivência 2D desenvolvido nativamente em **Go** com a engine gráfica **Ebitengine v2**. O jogo roda em uma janela desktop nativa, ambientado na **Idade Média / Fantasia Medieval**, com pixel art procedural em código, dispensando o carregamento de imagens externas.

---

## Estrutura do Projeto

```
├── cmd/
│   └── game/
│       └── main.go           # Ponto de entrada da aplicação desktop
├── internal/
│   ├── sprites/
│   │   └── sprites.go        # Gerador procedural de Pixel Art (Heróis, Inimigos e Drops)
│   ├── collision/
│   │   ├── aabb.go           # Detecção de colisões 2D (AABB / Bounding Boxes)
│   │   └── aabb_test.go      # Testes unitários de colisão
│   ├── ui/
│   │   ├── button.go         # Botões estilizados medievais com madeira, ferro forjado e rebites
│   │   └── bar.go            # Barras dinâmicas de progresso (HP, Mana)
│   ├── entities/
│   │   ├── player.go         # Jogador com Pixel Art, movimentação, stats e sistema de Buffs/Debuffs
│   │   ├── enemy.go          # 4 Arquétipos de inimigos (Comum, Corredor, Tanque e Atirador)
│   │   ├── drop.go           # Sistema de Drops (Loot Benigno e Maligno com animação flutuante)
│   │   ├── projectile.go     # Projéteis do jogador e cuspe ácido dos atiradores
│   │   └── slash.go          # Arco cortante do ataque Melee de espada do Guerreiro
│   ├── scenes/
│   │   ├── scene.go          # Interface Scene e SceneNavigator
│   │   ├── menu.go           # Menu Principal Medieval (muralhas, tochas com labaredas e brasão)
│   │   ├── options.go        # Configurações do Reino (Fullscreen on/off, ajuste de volume)
│   │   ├── class_select.go   # Seleção das 3 classes com prévia em Pixel Art
│   │   ├── gameplay.go       # Loop do combate, hordas, drops e Menu de Pausa com ESC
│   │   └── game_over.go      # Cripta da derrota com estatísticas e botões de reinício
│   └── game/
│       └── game.go           # Implementação de ebiten.Game (Update, Draw, Layout)
├── go.mod
└── go.sum
```

---

## Controles

- **ESC**: **Abre o Menu de Pausa** durante a batalha (permite retomar o combate, alternar tela cheia ou **voltar com segurança ao Menu Principal**).
- **W, A, S, D** ou **Setas do Teclado**: Movimentação do personagem.
- **Mouse (Cursor)**: Direção da mira dos ataques e interação com os menus.
- **Clique Esquerdo do Mouse** ou **Barra de Espaço**: Desferir ataque da classe.

---

## Ambientação Idade Média

- **Muralhas de Pedra Talhada:** Cenários e menus com blocos de rocha de fortaleza medieval com relevo e argamassa rústica.
- **Tochas Animadas:** Duas tochas em suportes de ferro forjado laterais com fogo crepitante em tempo real e fagulhas/brasas subindo.
- **Brasão Feudal e Espadas Cruzadas:** Escudo heráldico bipartido em carmesim e azul ultramar com espadas de cruzada entrecruzadas.
- **Placas e Botões de Carvalho Forjado:** Molduras duplas de ferro forjado com rebites dourados metálicos nos quatro cantos e ranhuras de madeira.

---

## Como Executar

```powershell
go run ./cmd/game
```

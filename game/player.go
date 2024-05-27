package game

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/ThreeDotsLabs/meteors/assets"
)

const (
	shootCooldown     = time.Millisecond * 30
	rotationPerSecond = math.Pi

	bulletSpawnOffset = 50.0
)

type Player struct {
	game *Game

	position      Vector
	rotation      float64
	sprite        *ebiten.Image
	lastAngle     Vector
	shootCooldown *Timer
	movSpeed      float64
	freezeMeteor  bool
}

func NewPlayer(game *Game) *Player {
	sprite := assets.PlayerSprite

	bounds := sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	pos := Vector{
		X: screenWidth/2 - halfW,
		Y: screenHeight/2 - halfH,
	}
	movSpeed := float64(7)
	return &Player{
		game:          game,
		position:      pos,
		rotation:      0,
		sprite:        sprite,
		shootCooldown: NewTimer(shootCooldown),
		movSpeed:      movSpeed,
	}
}

func (p *Player) Update() {
	speed := rotationPerSecond / float64(ebiten.TPS())

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.rotation -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.rotation += speed
	}

	if ebiten.IsKeyPressed(ebiten.KeyF) {
		p.freezeMeteor = !p.freezeMeteor
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		// Calculate the direction vector based on the current rotation
		dx := math.Sin(p.rotation) * p.movSpeed
		dy := -math.Cos(p.rotation) * p.movSpeed

		p.lastAngle.X = dx
		p.lastAngle.Y = dy
		// Update the player's position
		p.position.X += dx
		p.position.Y += dy
	} else {
		p.position.X += p.lastAngle.X / 10
		p.position.Y += p.lastAngle.Y / 10
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) {
		// Calculate the direction vector based on the current rotation
		dx := math.Sin(p.rotation) * p.movSpeed
		dy := -math.Cos(p.rotation) * p.movSpeed

		p.lastAngle.X = dx
		p.lastAngle.Y = dy

		// Update the player's position
		p.position.X -= dx
		p.position.Y -= dy

	} else {
		p.position.X += p.lastAngle.X / 10
		p.position.Y += p.lastAngle.Y / 10
	}

	p.shootCooldown.Update()
	if p.shootCooldown.IsReady() && ebiten.IsKeyPressed(ebiten.KeySpace) {
		p.shootCooldown.Reset()

		bounds := p.sprite.Bounds()
		halfW := float64(bounds.Dx()) / 2
		halfH := float64(bounds.Dy()) / 2

		spawnPos := Vector{
			p.position.X + halfW + math.Sin(p.rotation)*bulletSpawnOffset,
			p.position.Y + halfH + math.Cos(p.rotation)*-bulletSpawnOffset,
		}

		bullet := NewBullet(spawnPos, p.rotation)
		p.game.AddBullet(bullet)
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	bounds := p.sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(p.rotation)
	op.GeoM.Translate(halfW, halfH)

	op.GeoM.Translate(p.position.X, p.position.Y)

	screen.DrawImage(p.sprite, op)
}

func (p *Player) Collider() Rect {
	bounds := p.sprite.Bounds()

	return NewRect(
		p.position.X,
		p.position.Y,
		float64(bounds.Dx()),
		float64(bounds.Dy()),
	)
}

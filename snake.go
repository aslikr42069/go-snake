package main

import (
 _ "image/png"
 "math/rand"
 "time"

 "log"
 "github.com/hajimehoshi/ebiten/v2"
 "github.com/hajimehoshi/ebiten/v2/ebitenutil"
 )

type bodyNode struct{ /* Each element of the body says where it is located on the board */
 x, y int;
}

type bodyInfo struct { /* The body (and head) of the snake will be an array*/
 size, end int;        // end is equal to the index of the last node + 1
                       // The beginning of the body is always at index 0
}

func addNode(info bodyInfo, snakeBody []bodyNode, next_x int, next_y int) (bodyInfo, []bodyNode){
 tmp_Body := make([]bodyNode, info.size);
 tmp_Body[0].x = next_x;
 tmp_Body[0].y = next_y;
 for i := 0; i < info.end; i++ {
  tmp_Body[i+1].x = snakeBody[i].x;
  tmp_Body[i+1].y = snakeBody[i].y;
 }
 newInfo := bodyInfo{info.size, info.end + 1};
 return newInfo, tmp_Body;
}

func moveBody(info bodyInfo, snakeBody []bodyNode, next_x int, next_y int) []bodyNode {
 tmp_Body := make([]bodyNode, info.size);
 tmp_Body[0].x = next_x;
 tmp_Body[0].y = next_y;
 for i := 1; i < info.end; i++ {
  tmp_Body[i].x = snakeBody[i-1].x;
  tmp_Body[i].y = snakeBody[i-1].y;
 }
 return tmp_Body;
}

func generateApple(info bodyInfo, snakeBody []bodyNode, width int) (bodyNode) {
 appleCoords := bodyNode{rand.Intn(width), rand.Intn(width)};
 for i := 0; i < info.end; i++{
  if (appleCoords.x == snakeBody[i].x && appleCoords.y == snakeBody[i].y){
   return generateApple(info, snakeBody, width);
  }
 }
 return appleCoords;
}

/* end of snake transformation and translation code */


type Game struct{
 snakeInfo bodyInfo;
 snakeBody []bodyNode;
 apple     bodyNode;
}

var playing_field *ebiten.Image;
var snake_bod *ebiten.Image;
var apple_img *ebiten.Image;
var bg *ebiten.Image;
var game_over bool = false
var next_pos bodyNode = bodyNode{1, 0};


func init(){
 var err error;
 playing_field, _, err = ebitenutil.NewImageFromFile("images/playing_field.png");
 if err != nil {
  log.Fatal(err);
 }
 bg, _, err = ebitenutil.NewImageFromFile("images/playing_field.png");
 if err != nil {
  log.Fatal(err);
 }
 snake_bod, _, err = ebitenutil.NewImageFromFile("images/snake_body.png");
 if err != nil {
  log.Fatal(err);
 }
 apple_img, _, err = ebitenutil.NewImageFromFile("images/fruit.png");
 if err != nil {
  log.Fatal(err);
 }
 rand.Seed(time.Now().UnixNano());
}

func (g *Game) Update() error {
 if(game_over == true) {
  if(ebiten.IsKeyPressed(ebiten.KeySpace)){ 
   g.snakeInfo.end = 1;
   g.snakeBody[0].x = 10;
   g.snakeBody[0].y = 10;
   g.apple = generateApple(g.snakeInfo, g.snakeBody, 32);
   game_over = false;
  }
 } else {
 if(ebiten.IsKeyPressed(ebiten.KeyArrowDown)){
  var tmp_nextpos bodyNode = next_pos;
  next_pos.x = 0;
  next_pos.y = 1;
  for i := 0; i < g.snakeInfo.end; i++ {
   if(g.snakeBody[i].y ==  g.snakeBody[0].y + 1 && g.snakeBody[i].x ==  g.snakeBody[0].x){
    next_pos = tmp_nextpos;
   }
  }
 } else if(ebiten.IsKeyPressed(ebiten.KeyArrowUp)){
  var tmp_nextpos bodyNode = next_pos;
  next_pos.x = 0;
  next_pos.y = -1;
  for i := 0; i < g.snakeInfo.end; i++ {
   if(g.snakeBody[i].y ==  g.snakeBody[0].y - 1 && g.snakeBody[i].x ==  g.snakeBody[0].x){
    next_pos = tmp_nextpos;
   }
  }
 } else if(ebiten.IsKeyPressed(ebiten.KeyArrowLeft)){
  var tmp_nextpos bodyNode = next_pos;
  next_pos.x = -1;
  next_pos.y = 0;
  for i := 0; i < g.snakeInfo.end; i++ {
   if(g.snakeBody[i].y ==  g.snakeBody[0].y && g.snakeBody[i].x ==  g.snakeBody[0].x - 1){
    next_pos = tmp_nextpos;
   }
  }
 } else if(ebiten.IsKeyPressed(ebiten.KeyArrowRight)){
  var tmp_nextpos bodyNode = next_pos;
  next_pos.x = 1;
  next_pos.y = 0;
  for i := 0; i < g.snakeInfo.end; i++ {
   if(g.snakeBody[i].y ==  g.snakeBody[0].y && g.snakeBody[i].x ==  g.snakeBody[0].x + 1){
    next_pos = tmp_nextpos;
   }
  }
 }
 
 if(g.snakeBody[0].x == 1 && next_pos.x == -1){
  game_over = true;
 } else if(g.snakeBody[0].x == 32 && next_pos.x == 1){
  game_over = true;
 } else if(g.snakeBody[0].y == 32 && next_pos.y == 1){
  game_over = true;
 } else if(g.snakeBody[0].y == 1 && next_pos.y == -1){
  game_over = true;
 }

 for i := 0; i < g.snakeInfo.end; i++ {
  if((g.snakeBody[0].x + next_pos.x == g.snakeBody[i].x) && (g.snakeBody[0].y + next_pos.y == g.snakeBody[i].y)){
   game_over = true;
  }
 }

 if((g.snakeBody[0].x + next_pos.x == g.apple.x) && (g.snakeBody[0].y + next_pos.y == g.apple.y)){
  g.snakeInfo, g.snakeBody = addNode(g.snakeInfo, g.snakeBody, g.snakeBody[0].x + next_pos.x, g.snakeBody[0].y + next_pos.y);
  g.apple = generateApple(g.snakeInfo, g.snakeBody, 32);
 } else {
  g.snakeBody = moveBody(g.snakeInfo, g.snakeBody, g.snakeBody[0].x + next_pos.x, g.snakeBody[0].y + next_pos.y);
 }
 time.Sleep(100 * time.Millisecond);
 }
 return nil;
}

func (g *Game) Draw(screen *ebiten.Image){
 pos := &ebiten.DrawImageOptions{}
 playing_field.Clear();
 playing_field = ebiten.NewImageFromImage(bg);
 for i:= 0; i < g.snakeInfo.end; i++ {
  pos.GeoM.Translate(float64(g.snakeBody[i].x * 15), float64(g.snakeBody[i].y * 15));
  playing_field.DrawImage(snake_bod, pos);
  pos.GeoM.Translate(float64(g.snakeBody[i].x * -15), float64(g.snakeBody[i].y * -15));
 }
 pos.GeoM.Translate(float64(g.apple.x * 15), float64(g.apple.y * 15));
 playing_field.DrawImage(apple_img, pos);
 screen.DrawImage(playing_field, nil);
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int){
 return 480, 480;
}

func main(){
 const board_size = 32; // the size of the grid in which the snake exists will be 32x32 :))
 g := &Game{};
 g.snakeInfo.size = board_size * board_size;
 g.snakeBody = make([]bodyNode, g.snakeInfo.size);
 g.snakeInfo.end = 1;
 g.snakeBody[0].x = 10;
 g.snakeBody[0].y = 10;

 g.apple = generateApple(g.snakeInfo, g.snakeBody, board_size);

 ebiten.SetWindowSize(15*board_size, 15*board_size);
 ebiten.SetWindowTitle("Snake");
 if err := ebiten.RunGame(g); err != nil {
  log.Fatal(err);
 }
}


package main

// import (
// 	"fmt"
// 	"time"

// 	"fyne.io/fyne/v2/app"
// 	"fyne.io/fyne/v2/container"
// 	"fyne.io/fyne/v2/widget"
// )

// func main() {
// 	a := app.New()
// 	w := a.NewWindow("test")
// 	paused := false
// 	pauseCh := make(chan bool)
// 	resumeCh := make(chan bool)
// 	breakCh := make(chan bool)
// 	v1s := []string{"sss", "bbb", "sss", "bbb", "sss", "bbb"}
// 	v2s := []string{"111", "222", "111", "222", "111", "222"}

// 	w.SetContent(container.NewHBox(
// 		widget.NewButton("start", func() {
// 			go func() {
// 				j := 0
// 				for i := 0; i < len(v1s); {
// 					for j < len(v2s) {
// 						select {
// 						case <-pauseCh:
// 							paused = true
// 						case <-resumeCh:
// 							paused = false
// 						case <-breakCh:
// 							return
// 						default:
// 							if !paused {
// 								fmt.Println(v1s[i] + v2s[j])
// 								time.Sleep(time.Second * 1)
// 								j++
// 							}
// 						}
// 					}
// 					i++
// 					j = 0
// 				}
// 			}()
// 		}),
// 		widget.NewButton("pause", func() {
// 			if !paused {
// 				pauseCh <- true
// 			} else {
// 				resumeCh <- true
// 			}
// 		}),
// 		widget.NewButton("break", func() {
// 			breakCh <- true
// 		}),
// 	))

// 	w.ShowAndRun()
// }

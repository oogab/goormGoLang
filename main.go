package main

import "fmt"

type item struct {
	name string
	price int
	amount int
}

type buyer struct {
	point int
	shoppingBucket map[string]int
}

// 장바구니 목록은 사용자가 장바구니에 담은 물품의 이름(string)과 수량(int)을 맵형식으로 저장합니다.
// 따라서 생성자도 만드는 것이 좋습니다.
// 생성자 까먹었따!
func newBuyer() *buyer {
	d := buyer{}
	d.point = 1000000
	d.shoppingBucket = map[string]int{}
	return &d
}

func buying(item []item, byr *buyer, itemchoice int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			fmt.Println()
		}
	}()

	inputamount := 0 // 구매 수량

	fmt.Print("수량을 입력하시오 :")
	fmt.Scanln(&inputamount)
	fmt.Println() // 한 줄 띄어쓰기

	if inputamount <= 0 {
		panic("올바른 수량을 입력하세요.")
	}

	if byr.point < item[itemchoice-1].price * inputamount || item[itemchoice-1].amount < inputamount {
		panic("주문이 불가능합니다.")
	} else {
		for {
			buy := 0 // 살지 장바구니 담을지
			fmt.Println("1. 바로 주문\n2. 장바구니에 담기")
			fmt.Print("실행할 기능을 입력하시오 :")
			fmt.Scanln(&buy)
			fmt.Println()
	
			if buy == 1 { // 바로 주문
				item[itemchoice-1].amount -= inputamount
				byr.point -= item[itemchoice-1].price * inputamount
	
				fmt.Println("상품이 주문 접수되었습니다.")
				break
			} else if buy == 2 { // 장바구니에 담기
				checkbucket := false	// 중복 물품을 체크하기 위한 변수

				for itms := range byr.shoppingBucket { // 물품 체크
					if itms == item[itemchoice-1].name {
						checkbucket = true
					}
				}

				if checkbucket { // 장바구니에 중복되는 물품이 있을 때
					if byr.shoppingBucket[item[itemchoice-1].name] + inputamount > item[itemchoice-1].amount {
						fmt.Println("물품의 잔여 수량을 초과했습니다.")
						break
					}
					byr.shoppingBucket[item[itemchoice-1].name] += inputamount // 수량만 더함
				} else {	// 장바구니에 중복되는 물품이 없을 때
					byr.shoppingBucket[item[itemchoice-1].name] = inputamount	// 새로 품목 추가
				}

				fmt.Println("상품이 장바구니에 담겼습니다.")
				break // 구매 for문을 빠져나감
			} else {
				fmt.Println("잘못된 입력입니다. 다시 입력해주세요.")
			}
		}
	}
}

func main() {
	items := make([]item, 5) // 물품 목록
	buyer := newBuyer()			 // 구매자 정보(장바구니, 마일리지)

	items[0] = item{"텀블러", 10000, 30}
	items[1] = item{"롱패딩", 500000, 20}
	items[2] = item{"투미 백팩", 400000, 20}
	items[3] = item{"나이키 운동화", 150000, 50}
	items[4] = item{"빼빼로", 1200, 500}

	for {
		menu := 0 // 첫 메뉴

		fmt.Println("1. 구매")
		fmt.Println("2. 잔여 수량 확인")
		fmt.Println("3. 잔여 마일리지 확인")
		fmt.Println("4. 배송 상태 확인")
		fmt.Println("5. 장바구니 확인")
		fmt.Println("6. 프로그램 종료")
		fmt.Print("실행할 기능을 입력하시오 :")

		fmt.Scanln(&menu)
		fmt.Println()

		if menu == 1 { // 물건 구매
			for {
				itemchoice := 0
	
				for i := 0; i < 5; i++ {
					fmt.Printf("물품%d: %s, 가격: %d원, 잔여 수량: %d\n", i+1, items[i].name, items[i].price, items[i].amount)
				}
				fmt.Print("구매할 물품을 선택하세요 :")
				fmt.Scanln(&itemchoice)
				fmt.Println()
	
				if itemchoice == 1 {
					buying(items, buyer, 1)
					break
				} else if itemchoice == 2 {
					buying(items, buyer, 2)
					break
				} else if itemchoice == 3 {
					buying(items, buyer, 3)
					break
				} else if itemchoice == 4 {
					buying(items, buyer, 4)
					break
				} else if itemchoice == 5 {
					buying(items, buyer, 5)
					break
				} else {
					fmt.Println("잘못된 입력입니다. 다시 입력해주세요")
					fmt.Println()
				}
			}

			fmt.Print("엔터를 입력하면 메뉴 화면으로 돌아갑니다.")
			fmt.Scanln()
			fmt.Println()
		} else if menu == 2 { // 남은 수량 확인
			for i := 0; i < 5; i++ {
				fmt.Printf("%s, 잔여 수량: %d\n", items[i].name, items[i].amount)
			}
			fmt.Print("엔터를 입력하면 메뉴 화면으로 돌아갑니다.")
			fmt.Scanln()
			fmt.Println()
		} else if menu == 3 { // 잔여 마일리지 확인
			fmt.Printf("현재 잔여 마일리지는 %d점입니다.\n", buyer.point)
			fmt.Print("엔터를 입력하면 메뉴 화면으로 돌아갑니다.")
			fmt.Scanln()
			fmt.Println()
		} else if menu == 4 { // 배송 상태 확인

			fmt.Print("엔터를 입력하면 메뉴 화면으로 돌아갑니다.")
			fmt.Scanln()
			fmt.Println()
		} else if menu == 5 { // 장바구니 확인

			fmt.Print("엔터를 입력하면 메뉴 화면으로 돌아갑니다.")
			fmt.Scanln()
			fmt.Println()
		} else if menu == 6 { // 프로그램 종료
			fmt.Println("프로그램을 종료합니다.")
			return	// main함수 종료
		} else {
			fmt.Println("잘못된 입력입니다. 다시 입력해주세요.")
			fmt.Println()
		}
	}
}
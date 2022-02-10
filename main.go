package main

import (
	"fmt"
	"time"
)

type item struct {
	name string
	price int
	amount int
}

type buyer struct {
	point int
	shoppingBucket map[string]int
}

type delivery struct {
	status string
	onedelivery map[string]int	// 한번에 배송하는 물품의 뜻으로 생각
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

func newDelivery() delivery {
	d := delivery{}
	d.onedelivery = map[string]int{}
	return d
}

func buying(item []item, byr *buyer, itemchoice int, num *int, c chan bool, temp map[string]int) {
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
				if *num < 5 {
					item[itemchoice-1].amount -= inputamount
					byr.point -= item[itemchoice-1].price * inputamount
					temp[item[itemchoice-1].name] = inputamount

					c <- true
					*num++

					fmt.Println("상품이 주문 접수되었습니다.")
					fmt.Println()
					break
				} else {
					fmt.Println("배송 한도를 초과했습니다. 배송이 완료되면 주문하세요.")
					break
				}
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
						fmt.Println()
						break
					}
					byr.shoppingBucket[item[itemchoice-1].name] += inputamount // 수량만 더함
				} else {	// 장바구니에 중복되는 물품이 없을 때
					byr.shoppingBucket[item[itemchoice-1].name] = inputamount	// 새로 품목 추가
				}

				fmt.Println("상품이 장바구니에 담겼습니다.")
				fmt.Println()
				break // 구매 for문을 빠져나감
			} else {
				fmt.Println("잘못된 입력입니다. 다시 입력해주세요.")
				fmt.Println()
			}
		}
	}
}

func emptyBucket(byr *buyer) {
	if len(byr.shoppingBucket) == 0 {
		fmt.Println("장바구니가 비었습니다.")
	} else {
		for index, val := range byr.shoppingBucket {
			fmt.Printf("%s, 수량: %d\n", index, val)
		}
	}
	fmt.Println()
}

func requiredPoint(item []item, byr *buyer) (canbuy bool) {
	totalPoint := 0

	for index, val := range byr.shoppingBucket {
		for _, itms := range item {
			if index == itms.name {
				totalPoint += itms.price * val
			}
		}
	}

	fmt.Printf("필요 마일리지 : %d\n", totalPoint)
	fmt.Printf("보유 마일리지 : %d\n", byr.point)
	fmt.Println()

	if totalPoint > byr.point {
		fmt.Printf("마일리지가 %d점 부족합니다.\n", totalPoint - byr.point)
		return false
	}

	return true
}

func excessAmount(item []item, byr *buyer) (canbuy bool) {
	for index, val := range byr.shoppingBucket {
		for _, itms := range item {
			if index == itms.name {
				if itms.amount < val {	// 장바구니의 물품 총 개수가 판매하는 물품 개수보다 클 때
					fmt.Printf("%s, %d개 초과\n", itms.name, val - itms.amount)
					return false
				}
			}
		}
	}

	return true
}

func bucketBuying(item []item, byr *buyer, num *int, c chan bool, temp map[string]int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r, "\n")
		}
	}()

	if len(byr.shoppingBucket) == 0 {
		panic("주문 가능한 목록이 없습니다.")
	} else {
		if *num < 5 {
			for index, val := range byr.shoppingBucket {
				temp[index] = val	// 임시 저장

				for _, itms := range item {
					if itms.name == index {
						byr.point -= val * itms.price // 포인트 차감
						itms.amount -= val // 수량 차감
					}
				}
			}

			c <- true
			byr.shoppingBucket = map[string]int{} // 장바구니 초기화
			*num++
		} else {
			fmt.Println("배송 한도를 초과했습니다. 배송이 완료되면 주문하세요.")
		}
	}
}

func deliveryStatus(num *int, c chan bool, deliverylist []delivery, i int, temp *map[string]int) {
	for {
		if <-c {
			for index, val := range *temp {
				deliverylist[i].onedelivery[index] = val // 임시 저장한 데이터를 배송 상품에 저장함
			}

			*temp = map[string]int{}	// 임시 데이터 초기화
			
			deliverylist[i].status = "주문접수"
			time.Sleep(time.Second * 10)

			deliverylist[i].status = "배송중"
			time.Sleep(time.Second * 30)

			deliverylist[i].status = "배송완료"
			time.Sleep(time.Second * 10)

			deliverylist[i].status = ""
			*num--
			deliverylist[i].onedelivery = map[string]int{}	// 배송 리스트에서 물품 지우기
		}
	}
}

func main() {
	items := make([]item, 5) // 물품 목록
	buyer := newBuyer()			 // 구매자 정보(장바구니, 마일리지)
	numbuy := 0 // 주문한 개수
	deliverylist := make([]delivery, 5) // 배송 중인 상품 목록
	deliveryStart := make(chan bool, 5)
	tempdelivery := make(map[string]int) // 배달 물품 임시 저장

	for i := 0; i < 5; i++ {	// 배송 상품 객체 5개 생성
		deliverylist[i] = newDelivery()
	}

	for i := 0; i < 5; i++ {
		time.Sleep(time.Millisecond)	// 고루틴 순서대로 실행되도록 약간 딜레이
		go deliveryStatus(&numbuy, deliveryStart, deliverylist, i, &tempdelivery)
	}

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
					buying(items, buyer, 1, &numbuy, deliveryStart, tempdelivery)
					break
				} else if itemchoice == 2 {
					buying(items, buyer, 2, &numbuy, deliveryStart, tempdelivery)
					break
				} else if itemchoice == 3 {
					buying(items, buyer, 3, &numbuy, deliveryStart, tempdelivery)
					break
				} else if itemchoice == 4 {
					buying(items, buyer, 4, &numbuy, deliveryStart, tempdelivery)
					break
				} else if itemchoice == 5 {
					buying(items, buyer, 5, &numbuy, deliveryStart, tempdelivery)
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
			total := 0
			for i := 0; i < 5; i++ {
				total += len(deliverylist[i].onedelivery)
			}
			if total == 0 {
				fmt.Println("배송중인 상품이 없습니다.")
			} else {
				for i := 0; i < len(deliverylist); i++ {
					if len(deliverylist[i].onedelivery) != 0 {	// 배송중인 항목만 출력
						for index, val := range deliverylist[i].onedelivery {
							fmt.Printf("%s %d개/ ", index, val)
						}
						fmt.Printf("배송상황: %s\n", deliverylist[i].status)
					}
				}
			}


			fmt.Print("엔터를 입력하면 메뉴 화면으로 돌아갑니다.")
			fmt.Println()
			fmt.Scanln()
		} else if menu == 5 { // 장바구니 확인
			bucketmenu := 0

			for {
				emptyBucket(buyer)
				
				fmt.Println("1. 장바구니 상품 주문")
				fmt.Println("2. 장바구니 초기화")
				fmt.Println("3. 메뉴로 돌아가기")
				fmt.Print("실행할 기능을 입력하시오 :")
				fmt.Scanln(&bucketmenu)

				if bucketmenu == 1 {	// 장바구니 상품 주문하는 기능
					canbuy := requiredPoint(items, buyer)
					// 살수 있는지 없는지 확인하는 canbuy 선언 및 초기화
					canbuy = excessAmount(items, buyer)
					if canbuy { // 주문
						bucketBuying(items, buyer, &numbuy, deliveryStart, tempdelivery)
						fmt.Println("주문이 완료되었습니다.")
						break
					} else {
						fmt.Println("주문이 실패하였습니다.")
						break
					}
				} else if bucketmenu == 2 {
					buyer.shoppingBucket = map[string]int{} // 상품 초기화
					fmt.Println("장바구니를 초기화했습니다.")
					break
				} else if bucketmenu == 3 {
					fmt.Println()
					break
				} else {
					fmt.Println("잘못된 입력입니다. 다시 입력해주세요.")
				}
			}
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
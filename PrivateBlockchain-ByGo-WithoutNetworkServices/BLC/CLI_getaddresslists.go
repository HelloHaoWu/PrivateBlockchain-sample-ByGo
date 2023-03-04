package BLC

import "fmt"

// ↓打印所有钱包地址
func (cli *CLI) addressLists() {
	fmt.Println("打印所有钱包地址:")
	wallets, _ := NewWallets()
	for address, _ := range wallets.Wallets {
		fmt.Printf("User Address: %s\n", address)
	}
}

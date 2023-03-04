package BLC

func (cli *CLI) createWallet() { // 通过一个方法获取wallet对象
	wallets, _ := NewWallets()
	wallets.CreateNewWallet()
	//fmt.Println(len(wallets.Wallets)) // ← 其实保存是整个保存, 但Println显示的时候会节省内容

	// 把所有数据保存起来
	wallets.SaveWallets()
}

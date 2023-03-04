package BLC

func (cli *CLI) printchain() {
	blockchain := GetBlockchainObject()
	defer blockchain.DB.Close()
	blockchain.Printchain()
}

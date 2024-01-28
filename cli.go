package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type CLI struct {
	bc *Blockchain
}

func (cli *CLI) Run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	// -data
	addBlockData := addBlockCmd.String("data", "", "Block data")

	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		logError(err)
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		logError(err)
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage: ")
	fmt.Println(" addblock -data BLOCK_DATA -> add block to the blockchain")
	fmt.Println(" printchain -> print all the blocks of the blockchain")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) addBlock(data string) {
	cli.bc.AddBlocks(data)
	fmt.Println("Success!")
}

func (cli *CLI) printChain() {
	bci := cli.bc.Iterator()

	for {
		block := bci.Next()

		fmt.Println("Id: ", block.Id)
		fmt.Println("Data: ", string(block.Data))
		fmt.Println("Timestamp: ", block.Timestamp)
		fmt.Printf("Hash: %X\n", block.Hash)
		fmt.Printf("Previous Block Hash: %X\n", block.PrevBlockHash)
		pow := NewProofOfWork(block)
		isValid := pow.validate()
		fmt.Printf("PoW: %s\n", strconv.FormatBool(isValid))
		fmt.Println()
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
)

func parseFlags(args []string) *model {
	var (
		flags    = pflag.NewFlagSet(args[0], pflag.ContinueOnError)
		host     = flags.StringP("host", "h", "", "modbus host")
		port     = flags.Uint16P("port", "p", 0, "modbus host port")
		slaveID  = flags.Uint16P("slaveid", "s", 1, "modbus slave ID")
		regs     = flags.StringArrayP("register", "r", nil, "array of registers in format '12000,h,short' where 12000 is the address, h is register type and short is the data type. See notes below.")
		interval = flags.Uint16P("interval", "i", 0, "interval in seconds, will continuously read modbus registers after specified interval, or single iteration if not set")
	)

	if err := flags.Parse(args[1:]); err != nil {
		usage(flags)
		os.Exit(1)
	}

	if *host == "" || *port == 0 || len(*regs) == 0 {
		fmt.Printf("Host, port and at least one register are required.\n\n")
		usage(flags)
		os.Exit(1)
	}

	registers := make([]register, 0, len(*regs))
	for _, r := range *regs {
		tokens := strings.Split(r, ",")
		address, err := strconv.Atoi(tokens[0])
		if err != nil {
			fmt.Println("could not parse address", err)
			usage(flags)
			os.Exit(1)
		}
		rtype := tokens[1]
		if rtype != "h" && rtype != "i" {
			fmt.Println("only h or i supported for register type")
			usage(flags)
			os.Exit(1)
		}
		dtype := tokens[2]
		if dtype != "short" && dtype != "integer" && dtype != "float" {
			fmt.Println("only short, integer and float supported for data type")
			usage(flags)
			os.Exit(1)
		}
		registers = append(registers, register{address, rtype, dtype, 0})
	}

	return &model{
		server{host: *host, port: uint16(*port), slaveID: uint16(*slaveID)},
		registers,
		*interval,
	}
}

func usage(f *pflag.FlagSet) {
	fmt.Printf("Usage: %s [flags]\n\n", os.Args[0])
	f.PrintDefaults()
	fmt.Printf("\nWhen defining the modbus registers to read, use a comma separated string with the address, "+
		"the register type and the data type E.g. \n\n%s -h localhost -p 502 --register 12000,H,Short --register 13000,I,Integer\n\n"+
		"The first register is holding type (h) and the second is input type (i). "+
		"This will read from register 12000 as a short and from register 13000 & 13001 because the type is 32 bit integer.\n"+
		"", os.Args[0])
	os.Exit(1)
}

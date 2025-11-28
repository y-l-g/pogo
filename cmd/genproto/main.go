package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	ProtocolVersion = 1

	PktTypeData     = 0x00
	PktTypeError    = 0x01
	PktTypeFatal    = 0x02
	PktTypeHello    = 0x03
	PktTypeShm      = 0x04
	PktTypeShutdown = 0x09
)

func main() {
	if err := generatePHP(); err != nil {
		panic(err)
	}
	if err := generateC(); err != nil {
		panic(err)
	}
	fmt.Println("Protocol constants generated.")
}

func generatePHP() error {
	content := fmt.Sprintf(`<?php

declare(strict_types=1);

namespace Pogo\Runtime;

interface ProtocolConstants
{
    public const PROTOCOL_VERSION = %d;

    public const TYPE_DATA = 0x%02X;
    public const TYPE_ERROR = 0x%02X;
    public const TYPE_FATAL = 0x%02X;
    public const TYPE_HELLO = 0x%02X;
    public const TYPE_SHM = 0x%02X;
    public const TYPE_SHUTDOWN = 0x%02X;
}
`, ProtocolVersion, PktTypeData, PktTypeError, PktTypeFatal, PktTypeHello, PktTypeShm, PktTypeShutdown)

	path := filepath.Join("lib", "Runtime", "ProtocolConstants.php")
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func generateC() error {
	content := fmt.Sprintf(`#ifndef POGO_CONSTS_H
#define POGO_CONSTS_H

#define PROTOCOL_VERSION  %d

#define PKT_TYPE_DATA     0x%02X
#define PKT_TYPE_ERROR    0x%02X
#define PKT_TYPE_FATAL    0x%02X
#define PKT_TYPE_HELLO    0x%02X
#define PKT_TYPE_SHM      0x%02X
#define PKT_TYPE_SHUTDOWN 0x%02X

#endif
`, ProtocolVersion, PktTypeData, PktTypeError, PktTypeFatal, PktTypeHello, PktTypeShm, PktTypeShutdown)

	return os.WriteFile("pogo_consts.h", []byte(content), 0644)
}

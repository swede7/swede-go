package parser

import "me.weldnor/swede/core/common"

type Error struct {
	StartPosition common.Position
	EndPosition   common.Position
	Message       string
}

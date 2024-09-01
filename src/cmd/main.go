package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dotabuff/manta"
	"github.com/dotabuff/manta/dota"
)

func main() {
	// Create a new parser instance from a file. Alternatively see NewParser([]byte)
	f, err := os.Open("my_replay.dem")
	if err != nil {
		log.Fatalf("unable to open file: %s", err)
	}
	defer f.Close()

	p, err := manta.NewStreamParser(f)
	if err != nil {
		log.Fatalf("unable to create parser: %s", err)
	}

	p.Callbacks.OnCSVCMsg_CreateStringTable(func(m *dota.CSVCMsg_CreateStringTable) error {
		fmt.Println("creating stringtable", m.GetName())
		return nil
	})

	// Register a callback, this time for the OnCUserMessageSayText2 event.
	p.Callbacks.OnCUserMessageSayText2(func(m *dota.CUserMessageSayText2) error {
		log.Printf("%s said: %s\n", m.GetParam1(), m.GetParam2())
		return nil
	})

	playerName := func(e *manta.Entity) string {
		ni, ok := e.GetInt32("m_pEntity.m_nameStringableIndex")
		if !ok {
			return "[no m_pEntity.m_nameStringableIndex]"
		}
		name, ok := p.LookupStringByIndex("EntityNames", ni)
		if !ok {
			return fmt.Sprintf("[no stringtable for %d]", ni)
		}
		return name
	}

	p.Callbacks.OnCCitadelUserMsg_HeroKilled(func(m *dota.CCitadelUserMsg_HeroKilled) error {
		attacker := p.FindEntity(m.GetEntindexAttacker())
		victim := p.FindEntity(m.GetEntindexVictim())
		fmt.Printf("%s killed %s\n", playerName(attacker), playerName(victim))
		return nil
	})

	//p.Callbacks.OnCMsgFireBullets(func(m *dota.CMsgFireBullets) error {
	//	log.Printf("??? %d", m.GetShotNumber())
	//	return nil
	//})

	// Start parsing the replay!
	if err := p.Start(); err != nil {
		log.Printf("error! %s\n", err.Error())
	}

	log.Printf("Parse Complete!\n")
}

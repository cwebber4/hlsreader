#Copyright 2021 Chris Webber

//package hlsparser
package main

import (
    "os"
    "log"
    "encoding/binary"
    "errors"
    "io"
)
//run: go run . path\to\file.ts 2> out.txt
//view binary in vim: %!xxd -b

const TSP_SIZE_BYTES int = 188
const HEADER_SIZE_BYTES int = 4

const shiftSyncByte uint8 = 24
const shiftTransportErrorInd uint8 = 23
const shiftPayloadUnitStartInd uint8 = 22
const shiftTransportPriority uint8 = 21
const shiftPid uint8 = 8
const shiftTransportScramblingControl uint8 = 6
const shiftAdaptationFieldControl uint8 = 4

//01000111 010 0000000010001 00 01 0100
const tspMaskSyncByte uint32 = 0x47 << shiftSyncByte
const tspMaskTransportErrorInd uint32 = 0x1 << shiftTransportErrorInd
const tspMaskPayloadUnitStartInd uint32 = 0x1 << shiftPayloadUnitStartInd
const tspMaskTransportPriority uint32 = 0x1 << shiftTransportPriority
const tspMaskPid uint32 = 0x1FFF << shiftPid
const tspMaskTransportScramblingControl uint32 = 0x3 << shiftTransportScramblingControl
const tspMaskAdaptationFieldControl uint32 = 0x3 << shiftAdaptationFieldControl
const tspMaskContinuityCounter uint32 = 0xf

type TransportStreamPacket struct {
    //has a PesPacket
    
    //Transport stream packets are 188 bytes long.
    
    //sync_byte
    //transport_error_indicator
    //payload_unit_start_indicator
    //transport_priority
    //PID
    //transport_scrambling_control
    //adaptation_field_control
    //continuity_counter
    header uint32
    
    //adaptation field
    
    dataBytes []uint8 //can calculate size by 184 minus the number of bytes in the adaptation_field.
    
}

type ProgramAssociationTable struct {
    //pointer_field              8    uimsbf
}

//type PesPacket struct {
//    
//}


//Add option to save each decoded file to disk, save each raw file to disk, or decode + merge all files and save that to disk.
    
//each .ts file is made up of multiple transport stream packets. Transport stream packets are 188 bytes.

//01000111 01000000 00010001 00010100
//sync_byte
//

//make a file reader to read a ts file on disc, parse out the packets and print out what is found.
func main() {
    if (len(os.Args) < 2) {
        log.Println("Expected filename as argument")
        os.Exit(1)
    }
    
    tsPacket, err := parseTsFile(os.Args[1])
    if (err != nil) {
        log.Fatal(err)
    }
    
    log.Printf("Headers: %#X\n", tsPacket.header)
}

func parseTsFile(fileName string) (*TransportStreamPacket, error) {
    file, err := os.Open(fileName);
    if (err != nil) {
        return nil, err;
    }
    defer file.Close()
    
    var headerBuffer []byte = make([]byte, HEADER_SIZE_BYTES)
    var dataBuffer []byte = make([]byte, TSP_SIZE_BYTES - HEADER_SIZE_BYTES)
    var tsPacket TransportStreamPacket
    
    readCount, err := file.Read(headerBuffer)
    if (readCount < HEADER_SIZE_BYTES) {
        return nil, errors.New("Expected a 4 byte header. Read only " + string(readCount) + " bytes.")
    } else if (err != nil) {
        return nil, err
    }
    
    var header uint32 = binary.BigEndian.Uint32(headerBuffer)

    for header & tspMaskSyncByte == tspMaskSyncByte {
        log.Printf("Detected Transport Stream Packet")
        
        tsPacket.header = header
        
        //TODO: process rest of the packet.
        //transport_error_indicator
        log.Printf("Transport Error Indicator: %#X\n", header & tspMaskTransportErrorInd >> shiftTransportErrorInd)
        
        //PID
        log.Printf("PID: %#X\n", header & tspMaskPid >> shiftPid)
        
        //payload_unit_start_indicator
        log.Printf("Payload Unit Start Indicator: %#X\n", header & tspMaskPayloadUnitStartInd >> shiftPayloadUnitStartInd)
        
        //transport_scrambling_control
        log.Printf("Transport Scrambling Control: %#X\n", header & tspMaskTransportScramblingControl >> shiftTransportScramblingControl)
        
        
        //adaptation_field_control
        log.Printf("Adaptation Field Control: %#X\n", header & tspMaskAdaptationFieldControl >> shiftAdaptationFieldControl)
        
        //continuity_counter
        log.Printf("Continuity Counter: %#X\n", header & tspMaskContinuityCounter)
        
        
        //TODO: read data bytes
        //kludge
        readCount, err = file.Read(dataBuffer)
        if (err == io.EOF) {
            break
        } else if (readCount < TSP_SIZE_BYTES - HEADER_SIZE_BYTES) {
            return nil, errors.New("Expected data.")
        } else if (err != nil) {
            return nil, err
        }
        log.Printf("Read data\n")
        
        //clearBuffer(headerBuffer)
        readCount, err = file.Read(headerBuffer)
        if (err == io.EOF) {
            break
        } else if (readCount < HEADER_SIZE_BYTES) {
            return nil, errors.New("Expected a 4 byte header. Read only " + string(readCount) + " bytes.")
        } else if (err != nil) {
            return nil, err
        }
        
        header = binary.BigEndian.Uint32(headerBuffer)
    }
    
    //error handle. it wasn't a ts. malformed if bytes remaining.
    //if (header & tspMaskSyncByte != tspMaskSyncByte) {
    //    err := errors.New("The file is not a transport stream file.")
    //    return nil, err
    //}
    
    return &tsPacket, nil
}

func clearBuffer(buffer []byte) {
    for ii, _ := range buffer {
        buffer[ii] = 0x0
    }
}














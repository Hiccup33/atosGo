package main
 
import (
   "debug/dwarf"
   "debug/macho"
   "errors"
   "fmt"
   "io"
   "log"
   "os"
   "path"
)
 
var gCurrentFile   string
var gEntry   *dwarf.Entry
var gTargetAddress uint64
var gDwarfData    *dwarf.Data
 
func processCompileUnit(theReader *dwarf.Reader, depth int, theEntry *dwarf.Entry) {
 
   // Process the entry
   gCurrentFile = ...（略）...
   gEntry = ...（略）...
 
   if (theEntry.Children) {
      processChildren(theReader, depth+1, false)
   }
 
   gCurrentFile = ""
 
}
 
func processSubprogram(theReader *dwarf.Reader, depth int, theEntry *dwarf.Entry) {
 
   var lowAddr       uint64;
   var highAddr   uint64;
 
   // Get the state we need
   lowVal  := theEntry.Val(dwarf.AttrLowpc);
   //highVal := theEntry.Val(dwarf.AttrHighpc);
 
   if (lowVal != nil) {
      //fmt.Println("------")
      //fmt.Println(reflect.TypeOf(lowVal))
      //fmt.Println(lowVal)
      //fmt.Println(highVal)
      //fmt.Println("------end")
 
      lowAddr  = lowVal.(uint64)
      //highAddr = highVal
 
   }
 
   var high uint64
   var highOK bool
   highField := theEntry.AttrField(dwarf.AttrHighpc)
 
   if highField != nil {
      switch highField.Class {
      case dwarf.ClassAddress:
         high, highOK = highField.Val.(uint64)
      case dwarf.ClassConstant:
         off, ok := highField.Val.(int64)
         if ok {
            high = lowAddr + uint64(off)
            highOK = true
         }
      }
   }
 
   if highOK {
      highAddr = high
   }
 
   // Check for a match
   if (gTargetAddress >= lowAddr && gTargetAddress < highAddr) {
      name := theEntry.Val(dwarf.AttrName)
      //line := theEntry.Val(dwarf.AttrDeclLine)
 
 
      lineNumber, err := findLine(gDwarfData, gTargetAddress, gEntry)
      if  err != nil {
         log.Printf("lineReader findLine error")
      }
 
      //fmt.Printf("line %v \n",lineNumber)
      fmt.Printf("++++++++%v (%v:%v)\n", name, path.Base(gCurrentFile), lineNumber)
   }
 
   // Process the entry
   if (theEntry.Children) {
      processChildren(theReader, depth+1, false)
   }
}
 
// 找行号
func findLine(dwarfData *dwarf.Data, address uint64, entry *dwarf.Entry) (uint64,error) {
   var ErrUnknownLine = errors.New("ErrUnknownLine")
 
   var lineNumber    uint64
   lineNumber = 0
 
 
   var lineReader *dwarf.LineReader
   var theErr    error
   if lineReader, theErr = dwarfData.LineReader(entry); theErr != nil {
      log.Printf("findLine  lineReader error")
      return 0,ErrUnknownLine
   }
 
   // findPC
   var line2 dwarf.LineEntry
 
   lineNumber, err := ...（略）...
   if  err != nil {
      log.Printf("lineReader findLine error")
      return 0,ErrUnknownLine
   }
 
   return lineNumber,nil
 
}
 
func readNextEntry(theReader *dwarf.Reader) *dwarf.Entry {
 
   // Read the entry
   theEntry, theErr := theReader.Next();
   if (theErr != nil) {
      fmt.Printf("ERROR: %v\n", theErr.Error())
      theEntry = nil;
   }
 
   return(theEntry)
}
 
func processEntry(theReader *dwarf.Reader, depth int, theEntry *dwarf.Entry) {
 
   // Process the entry
   switch theEntry.Tag {
   case dwarf.TagCompileUnit: processCompileUnit(theReader, depth, theEntry)
   case dwarf.TagSubprogram:  processSubprogram( theReader, depth, theEntry)
   default:
      if (theEntry.Children) {
         processChildren(theReader, depth+1, true)
      }
   }
}
 
func processChildren(theReader *dwarf.Reader, depth int, canSkip bool) {
   // Process the children
   if (canSkip) {
      theReader.SkipChildren();
   } else {
      for {
         theChild := readNextEntry(theReader);
         if (theChild == nil || theChild.Tag == 0) {
            break;
         }
         processEntry(theReader, depth, theChild);
      }
   }
}
 
type FileHeader struct {
   Magic  uint32
   Cpu    macho.Cpu
   SubCpu uint32
   Type   uint32
   Ncmd   uint32
   Cmdsz  uint32
   Flags  uint32
}
 
 
type NormalFile struct {
   Magic  uint32
   FileHeader
}
 
func main() {
   var targetFile       *macho.FatFile
   var dwarfData        *dwarf.Data
   var file            *macho.File
   var theErr          error
   var theFatErr        error
   var pathMacho        string
   var runtimeAddress    uint64
   var loadAddress          uint64
   var relativeAddress       uint64
   var segmentAddress    uint64
 
   var entry     *dwarf.Entry
   //var err         error
   var fatErr     error
   var name         string
   var fileName      string
 
   fmt.Println("Hello, World!")
   name = ""
   fileName = ""
 
   fmt.Println(name,fileName)
 
   pathMacho = ...（略）...
   loadAddress = ...（略）...
   runtimeAddress = ...（略）...
 
   fmt.Println("Hello, %#", macho.Cpu386)
 
   //f, err := os.Open(pathMacho)
 
   file, theErr = macho.Open(pathMacho) 
   targetFile, theFatErr = macho.OpenFat(pathMacho)
   if (theFatErr == nil) {
 
      if (len(targetFile.Arches) < 2) {
         os.Exit(1)
      }
      targetFile := targetFile.Arches[1]
      segmentAddress = ...（略）...
      relativeAddress = ...（略）...
 
      dwarfData, theFatErr = targetFile.DWARF()
 
      r := dwarfData.Reader()
 
      if entry, fatErr = r.SeekPC(relativeAddress); fatErr != nil {
         log.Print("Not Found ...")
         return
      } else {
         log.Print("Found ...")
      }
      //name := entry.Val(dwarf.AttrName)
      targetFile.Close()
 
      gTargetAddress = relativeAddress
      var addr =  gTargetAddress
 
      lines, err := dwarfData.LineReader(entry)
      if err != nil {
         return
      }
      var lentry dwarf.LineEntry
      if err := lines.SeekPC(addr, &lentry); err != nil {
         return
      }
 
      // Try to find the function name.
   FindFatName:
      for entry, err := r.Next(); entry != nil && err == nil; entry, err = r.Next() {
         if entry.Tag == dwarf.TagSubprogram {
            ranges, err := dwarfData.Ranges(entry)
            if err != nil {
               return
            }
            for _, pcs := range ranges {
               if pcs[0] <= addr && addr < pcs[1] {
                  var ok bool
                  // TODO: AT_linkage_name, AT_MIPS_linkage_name.
                  name, ok = entry.Val(dwarf.AttrName).(string)
                  if ok {
                     break FindFatName
                  }
               }
            }
 
            //fmt.Printf("name %v \n", entry.Val(dwarf.LineEntry{Address: addr}))
         }
 
      }
 
      fileName = path.Base(lentry.File.Name)
 
      var lineReader *dwarf.LineReader
      if lineReader, err = dwarfData.LineReader(entry); err != nil {
         log.Printf("lineReader error")
         return
      }
 
      // findPC
      var line2 dwarf.LineEntry
 
      lineNumber, err := findPC(lineReader, relativeAddress, &line2)
      if  err != nil {
         log.Printf("lineReader seekPC error 2")
         return
      }
 
      fmt.Printf("---------%v (%v:%v)\n", name, fileName, lineNumber)
 
   }else if (theErr == nil) {
 
      segmentAddress = file.Segment("__TEXT").Addr
 
      // Calculate the target address
      relativeAddress = runtimeAddress - loadAddress + segmentAddress//0x100000000
      gTargetAddress  = relativeAddress
 
      dwarfData, theErr = file.DWARF()
      gDwarfData = dwarfData
      processChildren(dwarfData.Reader(), 0, false)
 
      file.Close()
 
   }
 
}
 
func findPC(r *dwarf.LineReader, pc uint64, entry *dwarf.LineEntry) (uint64,error) {
   var ErrUnknownPC = errors.New("ErrUnknownPC")
 
   var lineNumber    uint64
   lineNumber = 0
 
   if err := r.Next(entry); err != nil {
      return 0,err
   }
   if entry.Address > pc {
      // We're too far. Start at the beginning of the table.
      r.Reset()
      if err := r.Next(entry); err != nil {
         return 0,err
      }
      if entry.Address > pc {
         // The whole table starts after pc.
         r.Reset()
         return 0,ErrUnknownPC
      }
   }
 
   // Scan until we pass pc, then back up one.
   //var prev dwarf.LineEntry
 
   for {
 
      var next dwarf.LineEntry
      pos := r.Tell()
      if err := r.Next(&next); err != nil {
         if err == io.EOF {
            return 0,ErrUnknownPC
         }
         return 0,err
      }
 
      if next.Address > pc {
         log.Printf("line no 2: %+v", entry.Line)
 
         if entry.EndSequence {
            // pc is in a hole in the table.
            return 0,ErrUnknownPC
         }
         // entry is the desired entry. Back up the
         // cursor to "next" and return success.
         r.Seek(pos)
         //*entry = prev
         return lineNumber,nil
      } else {
         ...（略）...
      }
 
      *entry = next
   }
}

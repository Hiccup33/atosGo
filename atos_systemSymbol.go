func main() {
 var theFile *macho.File;
 var theErr error;
 var pathMacho string;
 var runtimeAddress uint64;
 var loadAddress uint64
 var relativeAddress uint64
  
 var bestDistance uint64
 var currentDistance uint64
 var segmentAddress uint64
  
 //var syms [] macho.Symbol
  
 fmt.Println("Hello, World!")
  
 pathMacho = "/Users/mac/Downloads/test/res/CoreFoundation";
 loadAddress = ...（略）...
 runtimeAddress = ...（略）...
  
  
 theFile, theErr = macho.Open(pathMacho);
 if (theErr != nil) {
  fmt.Println("Hello, World! error")
  os.Exit(1)
 }
  
 segmentAddress = ...（略）...
 relativeAddress = ...（略）...
 bestDistance = ...（略）...
  
 fmt.Printf("Symbol: %d" , segmentAddress)
  
 // 符号表
 for _, sym := range theFile.Symtab.Syms {
  currentDistance = ...（略）...
  if relativeAddress >= sym.Value && currentDistance <= bestDistance {
   bestDistance = currentDistance
   fmt.Println("Symbol: " + sym.Name)
  }
  
 }
  
 theFile.Close()
  
 }

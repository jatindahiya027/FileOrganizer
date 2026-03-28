package organizer

// ExtensionMap maps lowercase file extensions to their category folder name
var ExtensionMap = map[string]string{
	// Images
	".jpg": "Images", ".jpeg": "Images", ".png": "Images", ".gif": "Images",
	".bmp": "Images", ".tiff": "Images", ".tif": "Images", ".webp": "Images",
	".heic": "Images", ".heif": "Images", ".svg": "Images", ".ico": "Images",
	".raw": "Images", ".cr2": "Images", ".cr3": "Images", ".nef": "Images",
	".arw": "Images", ".orf": "Images", ".rw2": "Images", ".pef": "Images",
	".srw": "Images", ".dng": "Images", ".avif": "Images", ".jxl": "Images",
	".psd": "Images", ".xcf": "Images", ".ai": "Images", ".eps": "Images",
	".indd": "Images", ".psb": "Images",

	// Videos
	".mp4": "Videos", ".mkv": "Videos", ".avi": "Videos", ".mov": "Videos",
	".wmv": "Videos", ".flv": "Videos", ".webm": "Videos", ".m4v": "Videos",
	".mpg": "Videos", ".mpeg": "Videos", ".3gp": "Videos", ".3g2": "Videos",
	".ts": "Videos", ".mts": "Videos", ".m2ts": "Videos", ".vob": "Videos",
	".ogv": "Videos", ".f4v": "Videos", ".asf": "Videos", ".rm": "Videos",
	".rmvb": "Videos", ".divx": "Videos", ".xvid": "Videos", ".prproj": "Videos",
	".drp": "Videos", ".kdenlive": "Videos", ".h264": "Videos", ".h265": "Videos",
	".hevc": "Videos",

	// Audio
	".mp3": "Audio", ".wav": "Audio", ".flac": "Audio", ".aac": "Audio",
	".ogg": "Audio", ".wma": "Audio", ".m4a": "Audio", ".opus": "Audio",
	".aiff": "Audio", ".aif": "Audio", ".mid": "Audio", ".midi": "Audio",
	".ape": "Audio", ".wv": "Audio", ".mka": "Audio", ".ra": "Audio",
	".au": "Audio", ".snd": "Audio", ".dts": "Audio", ".ac3": "Audio",
	".amr": "Audio", ".gsm": "Audio", ".caf": "Audio", ".mxf": "Audio",
	".dsf": "Audio", ".dff": "Audio",

	// Documents
	".pdf": "Documents", ".doc": "Documents", ".docx": "Documents", ".odt": "Documents",
	".rtf": "Documents", ".txt": "Documents", ".md": "Documents", ".tex": "Documents",
	".pages": "Documents", ".wpd": "Documents", ".wps": "Documents", ".xps": "Documents",
	".abw": "Documents", ".sdw": "Documents",

	// Spreadsheets
	".xls": "Spreadsheets", ".xlsx": "Spreadsheets", ".ods": "Spreadsheets",
	".numbers": "Spreadsheets", ".csv": "Spreadsheets", ".tsv": "Spreadsheets",

	// Presentations
	".ppt": "Presentations", ".pptx": "Presentations", ".odp": "Presentations",
	".key": "Presentations",

	// Ebooks
	".epub": "Ebooks", ".mobi": "Ebooks", ".azw": "Ebooks", ".azw3": "Ebooks",
	".fb2": "Ebooks", ".lit": "Ebooks", ".pdb": "Ebooks", ".djvu": "Ebooks",
	".cbz": "Ebooks", ".cbr": "Ebooks", ".cb7": "Ebooks", ".cbt": "Ebooks",

	// 3D Files
	".obj": "3D_Files", ".fbx": "3D_Files", ".stl": "3D_Files", ".blend": "3D_Files",
	".dae": "3D_Files", ".3ds": "3D_Files", ".max": "3D_Files", ".ma": "3D_Files",
	".mb": "3D_Files", ".c4d": "3D_Files", ".lwo": "3D_Files", ".lws": "3D_Files",
	".ply": "3D_Files", ".gltf": "3D_Files", ".glb": "3D_Files", ".usd": "3D_Files",
	".usda": "3D_Files", ".usdc": "3D_Files", ".usdz": "3D_Files", ".abc": "3D_Files",
	".x3d": "3D_Files", ".vrml": "3D_Files", ".wrl": "3D_Files", ".stp": "3D_Files",
	".step": "3D_Files", ".iges": "3D_Files", ".igs": "3D_Files", ".f3d": "3D_Files",
	".3mf": "3D_Files", ".amf": "3D_Files", ".skp": "3D_Files", ".rvt": "3D_Files",
	".dwg": "3D_Files", ".dxf": "3D_Files", ".ifc": "3D_Files", ".ztl": "3D_Files",
	".zbrush": "3D_Files", ".spp": "3D_Files", ".sbs": "3D_Files", ".sbsar": "3D_Files",
	".zpr": "3D_Files",

	// Executables & Scripts
	".exe": "Executables", ".msi": "Executables", ".bat": "Executables", ".cmd": "Executables",
	".ps1": "Executables", ".sh": "Executables", ".bash": "Executables", ".zsh": "Executables",
	".fish": "Executables", ".app": "Executables", ".pkg": "Executables",
	".deb": "Executables", ".rpm": "Executables", ".apk": "Executables",
	".appimage": "Executables", ".run": "Executables", ".elf": "Executables",
	".com": "Executables", ".jar": "Executables",

	// Archives
	".zip": "Archives", ".rar": "Archives", ".7z": "Archives", ".tar": "Archives",
	".gz": "Archives", ".bz2": "Archives", ".xz": "Archives", ".tgz": "Archives",
	".tbz2": "Archives", ".lz": "Archives", ".lzma": "Archives", ".zst": "Archives",
	".cab": "Archives", ".lha": "Archives", ".arj": "Archives", ".ace": "Archives",
	".sitx": "Archives", ".z": "Archives",

	// Disk Images
	".iso": "Disk_Images", ".img": "Disk_Images", ".bin": "Disk_Images",
	".cue": "Disk_Images", ".nrg": "Disk_Images", ".mdf": "Disk_Images",
	".mds": "Disk_Images", ".ccd": "Disk_Images", ".toast": "Disk_Images",
	".vhd": "Disk_Images", ".vhdx": "Disk_Images", ".vmdk": "Disk_Images",
	".ova": "Disk_Images", ".ovf": "Disk_Images",

	// Code
	".html": "Code", ".htm": "Code", ".css": "Code", ".js": "Code",
	".jsx": "Code", ".tsx": "Code", ".json": "Code",
	".xml": "Code", ".yaml": "Code", ".yml": "Code", ".toml": "Code",
	".ini": "Code", ".cfg": "Code", ".conf": "Code", ".env": "Code",
	".c": "Code", ".cpp": "Code", ".cc": "Code", ".cxx": "Code",
	".h": "Code", ".hpp": "Code", ".cs": "Code", ".java": "Code",
	".kt": "Code", ".swift": "Code", ".m": "Code", ".mm": "Code",
	".php": "Code", ".sql": "Code", ".r": "Code", ".scala": "Code",
	".dart": "Code", ".hs": "Code", ".clj": "Code", ".ex": "Code",
	".exs": "Code", ".erl": "Code", ".go": "Code", ".rs": "Code",
	".py": "Code", ".rb": "Code", ".pl": "Code", ".lua": "Code",
	".wasm": "Code", ".v": "Code", ".sv": "Code", ".asm": "Code",
	".s": "Code", ".vim": "Code", ".el": "Code", ".lisp": "Code",
	".ml": "Code", ".fs": "Code", ".fsx": "Code",

	// Fonts
	".ttf": "Fonts", ".otf": "Fonts", ".woff": "Fonts", ".woff2": "Fonts",
	".eot": "Fonts", ".fon": "Fonts", ".fnt": "Fonts", ".pfb": "Fonts",
	".pfm": "Fonts", ".afm": "Fonts", ".bdf": "Fonts", ".pcf": "Fonts",

	// Database
	".db": "Database", ".sqlite": "Database", ".sqlite3": "Database",
	".mdb": "Database", ".accdb": "Database", ".dbf": "Database",
	".dump": "Database", ".bak": "Database", ".frm": "Database",
	".ibd": "Database", ".ldf": "Database", ".ndf": "Database",
}

// CategoryColor maps categories to hex colors for the UI
var CategoryColor = map[string]string{
	"Images":        "#60a5fa",
	"Videos":        "#f472b6",
	"Audio":         "#a78bfa",
	"Documents":     "#34d399",
	"Spreadsheets":  "#6ee7b7",
	"Presentations": "#fbbf24",
	"Ebooks":        "#fb923c",
	"3D_Files":      "#f87171",
	"Executables":   "#e879f9",
	"Archives":      "#94a3b8",
	"Disk_Images":   "#475569",
	"Code":          "#38bdf8",
	"Fonts":         "#c084fc",
	"Database":      "#4ade80",
	"Others":        "#64748b",
}

// GetCategory returns the category for a given extension (lowercase)
func GetCategory(ext string) string {
	if cat, ok := ExtensionMap[ext]; ok {
		return cat
	}
	return "Others"
}

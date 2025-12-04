import { useEffect, useState } from "react";
import { FolderIcon } from "lucide-react";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Area } from "@/types/proto/api/v1/area_service";
import { Folder } from "@/types/proto/api/v1/folder_service";

interface Props {
  selectedFolder?: string;
  selectedArea?: string;
  onFolderChange?: (folderName: string | undefined) => void;
  onAreaChange?: (areaName: string | undefined) => void;
  className?: string;
}

const FolderSelector = ({ selectedFolder, selectedArea, onFolderChange, onAreaChange, className }: Props) => {
  const [areas, setAreas] = useState<Area[]>([]);
  const [folders, setFolders] = useState<Folder[]>([]);
  const [currentArea, setCurrentArea] = useState<string | undefined>(selectedArea);
  const [currentFolder, setCurrentFolder] = useState<string | undefined>(selectedFolder);

  useEffect(() => {
    fetchAreas();
    fetchFolders();
  }, []);

  const fetchAreas = async () => {
    try {
      const response = await fetch("/api/v1/areas");
      if (!response.ok) return;
      const data = await response.json();
      setAreas(data.areas || []);
    } catch (error) {
      console.error("Failed to fetch areas:", error);
    }
  };

  const fetchFolders = async () => {
    try {
      const response = await fetch("/api/v1/folders");
      if (!response.ok) return;
      const data = await response.json();
      setFolders(data.folders || []);
    } catch (error) {
      console.error("Failed to fetch folders:", error);
    }
  };

  const handleAreaChange = (areaName: string) => {
    if (areaName === "none") {
      setCurrentArea(undefined);
      setCurrentFolder(undefined);
      onAreaChange?.(undefined);
      onFolderChange?.(undefined);
    } else {
      setCurrentArea(areaName);
      setCurrentFolder(undefined);
      onAreaChange?.(areaName);
      onFolderChange?.(undefined);
    }
  };

  const handleFolderChange = (folderName: string) => {
    if (folderName === "none") {
      setCurrentFolder(undefined);
      onFolderChange?.(undefined);
    } else {
      const folder = folders.find((f) => f.name === folderName);
      if (folder) {
        setCurrentFolder(folderName);
        setCurrentArea(folder.area);
        onFolderChange?.(folderName);
        onAreaChange?.(folder.area);
      }
    }
  };

  const availableFolders = currentArea ? folders.filter((f) => f.area === currentArea) : [];

  return (
    <div className={className}>
      <div className="flex items-center gap-2">
        <FolderIcon className="w-4 h-4 text-gray-500" />
        <Select value={currentArea || "none"} onValueChange={handleAreaChange}>
          <SelectTrigger className="w-[180px] h-8">
            <SelectValue placeholder="Select area" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="none">No area</SelectItem>
            {areas.map((area) => (
              <SelectItem key={area.name} value={area.name}>
                {area.displayName}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>

        {currentArea && (
          <Select value={currentFolder || "none"} onValueChange={handleFolderChange}>
            <SelectTrigger className="w-[180px] h-8">
              <SelectValue placeholder="Select folder" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="none">No folder</SelectItem>
              {availableFolders.map((folder) => (
                <SelectItem key={folder.name} value={folder.name}>
                  {folder.displayName}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        )}
      </div>
    </div>
  );
};

export default FolderSelector;

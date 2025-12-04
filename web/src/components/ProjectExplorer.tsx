import { useEffect, useState } from "react";
import { toast } from "react-hot-toast";
import { ChevronDown, ChevronRight, FolderIcon, FolderPlus, Plus, LayoutGrid } from "lucide-react";
import { Area } from "@/types/proto/api/v1/area_service";
import { Folder } from "@/types/proto/api/v1/folder_service";
import CreateAreaDialog from "./CreateAreaDialog";
import CreateFolderDialog from "./CreateFolderDialog";
import { Button } from "@/components/ui/button";

const ProjectExplorer = () => {
  const [areas, setAreas] = useState<Area[]>([]);
  const [folders, setFolders] = useState<Folder[]>([]);
  const [expandedAreas, setExpandedAreas] = useState<Set<string>>(new Set());
  const [showCreateArea, setShowCreateArea] = useState(false);
  const [showCreateFolder, setShowCreateFolder] = useState(false);
  const [selectedAreaId, setSelectedAreaId] = useState<string>();

  useEffect(() => {
    fetchAreas();
    fetchFolders();
  }, []);

  const fetchAreas = async () => {
    try {
      const response = await fetch("/api/v1/areas", {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      });
      if (!response.ok) throw new Error("Failed to fetch areas");
      const data = await response.json();
      setAreas(data.areas || []);
    } catch (error) {
      console.error("Failed to fetch areas:", error);
      toast.error("Failed to load areas");
    }
  };

  const fetchFolders = async () => {
    try {
      const response = await fetch("/api/v1/folders", {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      });
      if (!response.ok) throw new Error("Failed to fetch folders");
      const data = await response.json();
      setFolders(data.folders || []);
    } catch (error) {
      console.error("Failed to fetch folders:", error);
      toast.error("Failed to load folders");
    }
  };

  const handleCreateArea = async (area: Partial<Area>) => {
    try {
      const response = await fetch("/api/v1/areas", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ area }),
      });
      if (!response.ok) throw new Error("Failed to create area");
      await fetchAreas();
    } catch (error) {
      console.error("Failed to create area:", error);
      throw error;
    }
  };

  const handleCreateFolder = async (folder: Partial<Folder>) => {
    try {
      const response = await fetch("/api/v1/folders", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ folder }),
      });
      if (!response.ok) throw new Error("Failed to create folder");
      await fetchFolders();
    } catch (error) {
      console.error("Failed to create folder:", error);
      throw error;
    }
  };

  const toggleArea = (areaName: string) => {
    const newExpanded = new Set(expandedAreas);
    if (newExpanded.has(areaName)) {
      newExpanded.delete(areaName);
    } else {
      newExpanded.add(areaName);
    }
    setExpandedAreas(newExpanded);
  };

  const getFoldersForArea = (areaName: string) => {
    return folders.filter((folder) => folder.area === areaName);
  };

  const handleOpenCreateFolder = (areaName: string) => {
    setSelectedAreaId(areaName);
    setShowCreateFolder(true);
  };

  return (
    <div className="w-full">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-lg font-semibold flex items-center gap-2">
          <LayoutGrid className="w-5 h-5" />
          Projects
        </h2>
        <Button size="sm" variant="outline" onClick={() => setShowCreateArea(true)}>
          <Plus className="w-4 h-4 mr-1" />
          New Area
        </Button>
      </div>

      <div className="space-y-1">
        {areas.length === 0 ? (
          <div className="text-sm text-gray-500 py-4 text-center">
            No areas yet. Create your first area to get started.
          </div>
        ) : (
          areas.map((area) => {
            const isExpanded = expandedAreas.has(area.name);
            const areaFolders = getFoldersForArea(area.name);

            return (
              <div key={area.name} className="space-y-0.5">
                <div className="flex items-center gap-1 group hover:bg-gray-100 dark:hover:bg-gray-800 rounded px-2 py-1.5">
                  <button
                    onClick={() => toggleArea(area.name)}
                    className="flex items-center gap-1 flex-1 text-sm font-medium"
                  >
                    {isExpanded ? (
                      <ChevronDown className="w-4 h-4 text-gray-500" />
                    ) : (
                      <ChevronRight className="w-4 h-4 text-gray-500" />
                    )}
                    <LayoutGrid className="w-4 h-4 text-blue-500" />
                    <span>{area.displayName}</span>
                  </button>
                  <button
                    onClick={() => handleOpenCreateFolder(area.name)}
                    className="opacity-0 group-hover:opacity-100 p-1 hover:bg-gray-200 dark:hover:bg-gray-700 rounded"
                    title="Create folder"
                  >
                    <FolderPlus className="w-4 h-4 text-gray-600 dark:text-gray-400" />
                  </button>
                </div>

                {isExpanded && (
                  <div className="ml-6 space-y-0.5">
                    {areaFolders.length === 0 ? (
                      <div className="text-xs text-gray-500 py-2 px-2">No folders</div>
                    ) : (
                      areaFolders.map((folder) => (
                        <div
                          key={folder.name}
                          className="flex items-center gap-2 hover:bg-gray-100 dark:hover:bg-gray-800 rounded px-2 py-1.5 cursor-pointer"
                        >
                          <FolderIcon className="w-4 h-4 text-yellow-500" />
                          <span className="text-sm">{folder.displayName}</span>
                        </div>
                      ))
                    )}
                  </div>
                )}
              </div>
            );
          })
        )}
      </div>

      <CreateAreaDialog open={showCreateArea} onClose={() => setShowCreateArea(false)} onConfirm={handleCreateArea} />

      <CreateFolderDialog
        open={showCreateFolder}
        areas={areas}
        selectedAreaId={selectedAreaId}
        onClose={() => {
          setShowCreateFolder(false);
          setSelectedAreaId(undefined);
        }}
        onConfirm={handleCreateFolder}
      />
    </div>
  );
};

export default ProjectExplorer;

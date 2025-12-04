import { useState } from "react";
import { toast } from "react-hot-toast";
import { useTranslation } from "react-i18next";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Folder } from "@/types/proto/api/v1/folder_service";
import { Area } from "@/types/proto/api/v1/area_service";

interface Props {
  open: boolean;
  areas: Area[];
  selectedAreaId?: string;
  onClose: () => void;
  onConfirm: (folder: Partial<Folder>) => Promise<void>;
}

const CreateFolderDialog = ({ open, areas, selectedAreaId, onClose, onConfirm }: Props) => {
  const { t } = useTranslation();
  const [displayName, setDisplayName] = useState("");
  const [description, setDescription] = useState("");
  const [areaName, setAreaName] = useState(selectedAreaId || "");
  const [isLoading, setIsLoading] = useState(false);

  const handleCreate = async () => {
    if (!displayName.trim()) {
      toast.error("Folder name is required");
      return;
    }
    if (!areaName) {
      toast.error("Please select an area");
      return;
    }

    setIsLoading(true);
    try {
      await onConfirm({
        displayName: displayName.trim(),
        description: description.trim(),
        area: areaName,
      });
      setDisplayName("");
      setDescription("");
      setAreaName(selectedAreaId || "");
      onClose();
      toast.success("Folder created successfully");
    } catch (error) {
      console.error("Failed to create folder:", error);
      toast.error("Failed to create folder");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Create New Folder</DialogTitle>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          <div className="grid gap-2">
            <Label htmlFor="area">Area</Label>
            <Select value={areaName} onValueChange={setAreaName}>
              <SelectTrigger>
                <SelectValue placeholder="Select an area" />
              </SelectTrigger>
              <SelectContent>
                {areas.map((area) => (
                  <SelectItem key={area.name} value={area.name}>
                    {area.displayName}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <div className="grid gap-2">
            <Label htmlFor="name">Name</Label>
            <Input
              id="name"
              value={displayName}
              onChange={(e) => setDisplayName(e.target.value)}
              placeholder="e.g., Equipment, Self-hosted apps"
            />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="description">Description</Label>
            <Textarea
              id="description"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="Optional description"
              rows={3}
            />
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={onClose} disabled={isLoading}>
            Cancel
          </Button>
          <Button onClick={handleCreate} disabled={isLoading}>
            {isLoading ? "Creating..." : "Create"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};

export default CreateFolderDialog;

import { useState } from "react";
import { toast } from "react-hot-toast";
import { useTranslation } from "react-i18next";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { Area } from "@/types/proto/api/v1/area_service";

interface Props {
  open: boolean;
  onClose: () => void;
  onConfirm: (area: Partial<Area>) => Promise<void>;
}

const CreateAreaDialog = ({ open, onClose, onConfirm }: Props) => {
  const { t } = useTranslation();
  const [displayName, setDisplayName] = useState("");
  const [description, setDescription] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const handleCreate = async () => {
    if (!displayName.trim()) {
      toast.error("Area name is required");
      return;
    }

    setIsLoading(true);
    try {
      await onConfirm({
        displayName: displayName.trim(),
        description: description.trim(),
      });
      setDisplayName("");
      setDescription("");
      onClose();
      toast.success("Area created successfully");
    } catch (error) {
      console.error("Failed to create area:", error);
      toast.error("Failed to create area");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Create New Area</DialogTitle>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          <div className="grid gap-2">
            <Label htmlFor="name">Name</Label>
            <Input
              id="name"
              value={displayName}
              onChange={(e) => setDisplayName(e.target.value)}
              placeholder="e.g., Homelabs, Work, Personal"
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

export default CreateAreaDialog;

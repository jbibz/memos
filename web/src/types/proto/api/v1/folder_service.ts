export interface Folder {
  name: string;
  uid: string;
  createTime?: Date;
  updateTime?: Date;
  displayName: string;
  description: string;
  area: string;
  parent?: string;
}

export interface CreateFolderRequest {
  folder: Folder;
}

export interface ListFoldersRequest {
  pageSize?: number;
  pageToken?: string;
  filter?: string;
}

export interface ListFoldersResponse {
  folders: Folder[];
  nextPageToken?: string;
}

export interface GetFolderRequest {
  name: string;
}

export interface UpdateFolderRequest {
  folder: Folder;
  updateMask?: string[];
}

export interface DeleteFolderRequest {
  name: string;
}

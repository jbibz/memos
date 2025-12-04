export interface Area {
  name: string;
  uid: string;
  createTime?: Date;
  updateTime?: Date;
  displayName: string;
  description: string;
  parent?: string;
}

export interface CreateAreaRequest {
  area: Area;
}

export interface ListAreasRequest {
  pageSize?: number;
  pageToken?: string;
  filter?: string;
}

export interface ListAreasResponse {
  areas: Area[];
  nextPageToken?: string;
}

export interface GetAreaRequest {
  name: string;
}

export interface UpdateAreaRequest {
  area: Area;
  updateMask?: string[];
}

export interface DeleteAreaRequest {
  name: string;
}

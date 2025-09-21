export interface UploadResponse {
  files: {
    id: string;
    name: string;
    url: string;
  }[];
}

export interface DeleteRequest {
  fileIds: string[];
}
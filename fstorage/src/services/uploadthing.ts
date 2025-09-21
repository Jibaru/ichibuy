import { UTApi } from 'uploadthing/server';

const utapi = new UTApi({
  token: process.env.UPLOADTHING_TOKEN,
});

export class UploadThingService {
  async uploadFiles(files: { name: string; data: Buffer }[]): Promise<{ id: string; name: string; url: string }[]> {
    try {
      const filesData = files.map(file => new File([new Uint8Array(file.data)], file.name));
      const response = await utapi.uploadFiles(filesData);

      return response.map((result, index) => {
        if (result.error) {
          throw new Error(`Failed to upload ${files[index].name}: ${result.error.message}`);
        }
        return {
          id: result.data.key,
          name: files[index].name,
          url: result.data.ufsUrl
        };
      });
    } catch (error) {
      console.error('Upload error:', error);
      throw new Error('Failed to upload files');
    }
  }

  async deleteFiles(fileIds: string[]): Promise<void> {
    try {
      await utapi.deleteFiles(fileIds);
    } catch (error) {
      console.error('Delete error:', error);
      throw new Error('Failed to delete files');
    }
  }
}
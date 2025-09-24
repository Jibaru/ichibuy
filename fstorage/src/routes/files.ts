import { Router, Response } from 'express';
import multer from 'multer';
import { UploadThingService } from '../services/uploadthing.js';
import { DeleteRequest } from '../types/index.js';
import { JWTAuthMiddleware, AuthenticatedRequest } from '../middlewares/jwt_auth.js';

const router = Router();
const uploadService = new UploadThingService();
const jwtAuth = new JWTAuthMiddleware(process.env.AUTH_BASE_URL);

const upload = multer({
  storage: multer.memoryStorage(),
  limits: {
    fileSize: 50 * 1024 * 1024, // 50MB limit
  },
});

router.post('/upload', jwtAuth.validateToken(), upload.array('files') as any, async (req: AuthenticatedRequest, res: Response) => {
  try {
    const files = req.files as Express.Multer.File[];
    const { domain } = req.body;

    if (!files || files.length === 0) {
      return res.status(400).json({ error: 'At least one file is required' });
    }

    if (!domain) {
      return res.status(400).json({ error: 'Domain is required' });
    }

    const filesWithBuffers = files.map(file => ({
      name: file.originalname,
      data: file.buffer
    }));

    const uploadedFiles = await uploadService.uploadFiles(filesWithBuffers);

    res.json({ files: uploadedFiles });
  } catch (error) {
    console.error('Upload endpoint error:', error);
    res.status(500).json({ error: 'Internal server error' });
  }
});

router.delete('/batch', jwtAuth.validateToken(), async (req: AuthenticatedRequest, res: Response) => {
  try {
    const { fileIds }: DeleteRequest = req.body;

    if (!fileIds || !Array.isArray(fileIds) || fileIds.length === 0) {
      return res.status(400).json({ error: 'File IDs array is required and cannot be empty' });
    }

    await uploadService.deleteFiles(fileIds);

    res.status(204).send();
  } catch (error) {
    console.error('Delete endpoint error:', error);
    res.status(500).json({ error: 'Internal server error' });
  }
});

export { router as filesRouter };
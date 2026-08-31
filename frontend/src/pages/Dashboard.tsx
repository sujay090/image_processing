import React, { useState, useEffect, useRef } from 'react';
import { Link } from 'react-router-dom';
import { Upload, Image as ImageIcon, Plus } from 'lucide-react';
import { apiFetch } from '../lib/api';
import { Button } from '../components/Button';
import { Card, CardContent } from '../components/Card';
import { Loader } from '../components/Loader';

interface ImageRecord {
  id: string;
  original_filename: string;
  format: string;
  size: number;
  width: number;
  height: number;
  status: string;
  created_at: string;
}

export function Dashboard() {
  const [images, setImages] = useState<ImageRecord[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isUploading, setIsUploading] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const fetchImages = async () => {
    try {
      setIsLoading(true);
      const data = await apiFetch('/images?limit=50');
      setImages(data || []);
    } catch (error) {
      console.error('Failed to fetch images', error);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchImages();
  }, []);

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    if (!e.target.files || e.target.files.length === 0) return;
    
    const file = e.target.files[0];
    const formData = new FormData();
    formData.append('image', file);

    try {
      setIsUploading(true);
      await apiFetch('/images', {
        method: 'POST',
        data: formData,
      });
      // Refresh list after upload
      await fetchImages();
    } catch (error) {
      console.error('Upload failed', error);
      alert('Failed to upload image. Please try again.');
    } finally {
      setIsUploading(false);
      if (fileInputRef.current) {
        fileInputRef.current.value = '';
      }
    }
  };

  return (
    <div className="container mx-auto px-4 pb-12">
      <div className="flex flex-col md:flex-row justify-between items-center mb-8 gap-4">
        <div>
          <h1 className="text-3xl font-bold mb-2">Dashboard</h1>
          <p className="text-gray-400">Manage and transform your uploaded images</p>
        </div>
        
        <div className="relative">
          <input 
            type="file" 
            ref={fileInputRef} 
            onChange={handleFileChange} 
            className="hidden" 
            accept="image/jpeg, image/png, image/webp, image/tiff, image/gif, image/avif"
          />
          <Button 
            onClick={() => fileInputRef.current?.click()} 
            isLoading={isUploading}
            className="whitespace-nowrap"
          >
            {!isUploading && <Upload className="w-4 h-4 mr-2" />}
            {isUploading ? 'Uploading...' : 'Upload Image'}
          </Button>
        </div>
      </div>

      {isLoading ? (
        <Loader fullPage={false} size="lg" className="mt-20" />
      ) : images.length === 0 ? (
        <Card className="mt-12 text-center p-12 border-dashed border-2 border-glass-border bg-glass/30">
          <CardContent className="flex flex-col items-center justify-center">
            <div className="p-4 bg-accent-primary/10 rounded-full mb-4">
              <ImageIcon className="w-12 h-12 text-accent-primary" />
            </div>
            <h3 className="text-xl font-semibold mb-2">No images yet</h3>
            <p className="text-gray-400 mb-6 max-w-sm">
              Upload your first image to start processing, resizing, and applying filters.
            </p>
            <Button onClick={() => fileInputRef.current?.click()} variant="secondary">
              <Plus className="w-4 h-4 mr-2" /> Select Image
            </Button>
          </CardContent>
        </Card>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          {images.map((img) => (
            <Link key={img.id} to={`/images/${img.id}`} className="group relative block animate-fade-in">
              <Card className="h-full transition-transform duration-300 group-hover:-translate-y-2 group-hover:shadow-[0_10px_40px_rgba(99,102,241,0.2)]">
                <div className="aspect-square bg-bg-tertiary flex items-center justify-center border-b border-glass-border p-8 relative overflow-hidden">
                  <div className="absolute inset-0 bg-gradient-to-t from-bg-primary/80 to-transparent opacity-0 group-hover:opacity-100 transition-opacity z-10"></div>
                  <ImageIcon className="w-16 h-16 text-gray-500 group-hover:scale-110 transition-transform duration-500" />
                  <div className="absolute bottom-4 left-4 right-4 z-20 opacity-0 group-hover:opacity-100 transition-opacity translate-y-4 group-hover:translate-y-0 duration-300">
                    <span className="bg-accent-primary text-white text-xs px-2 py-1 rounded-full font-medium">View Details</span>
                  </div>
                </div>
                <CardContent className="p-4">
                  <p className="font-semibold text-white truncate" title={img.original_filename}>
                    {img.original_filename}
                  </p>
                  <div className="flex justify-between items-center mt-2 text-xs text-gray-400">
                    <span className="uppercase px-2 py-0.5 bg-white/10 rounded">{img.format}</span>
                    <span>{(img.size / 1024).toFixed(1)} KB</span>
                  </div>
                </CardContent>
              </Card>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}

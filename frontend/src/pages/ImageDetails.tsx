import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { ArrowLeft, Wand2, RefreshCw, CheckCircle2, AlertCircle } from 'lucide-react';
import { apiFetch } from '../lib/api';
import { Button } from '../components/Button';
import { Input } from '../components/Input';
import { Card, CardContent, CardHeader } from '../components/Card';
import { Loader } from '../components/Loader';

export function ImageDetails() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [image, setImage] = useState<any>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isTransforming, setIsTransforming] = useState(false);
  const [transformStatus, setTransformStatus] = useState<any>(null);

  // Transformation Form State
  const [format, setFormat] = useState('jpeg');
  const [quality, setQuality] = useState('80');
  const [resizeW, setResizeW] = useState('');
  const [resizeH, setResizeH] = useState('');
  const [grayscale, setGrayscale] = useState(false);
  const [flip, setFlip] = useState(false);

  useEffect(() => {
    fetchImage();
  }, [id]);

  const fetchImage = async () => {
    try {
      setIsLoading(true);
      const data = await apiFetch(`/images/${id}`);
      setImage(data);
    } catch (error) {
      console.error(error);
      navigate('/');
    } finally {
      setIsLoading(false);
    }
  };

  const handleTransform = async (e: React.FormEvent) => {
    e.preventDefault();
    
    const transformations: any = {
      format: format,
      quality: parseInt(quality, 10),
    };

    if (resizeW && resizeH) {
      transformations.resize = { width: parseInt(resizeW, 10), height: parseInt(resizeH, 10) };
    }
    if (grayscale) {
      transformations.filters = { grayscale: true };
    }
    if (flip) {
      transformations.flip = true;
    }

    try {
      setIsTransforming(true);
      const res = await apiFetch(`/images/${id}/transform`, {
        method: 'POST',
        data: { transformations },
      });
      
      if (res.transformation_id) {
        pollStatus(res.transformation_id);
      }
    } catch (error) {
      console.error(error);
      alert('Failed to start transformation');
      setIsTransforming(false);
    }
  };

  const pollStatus = async (tfId: string) => {
    try {
      const data = await apiFetch(`/transformations/${tfId}`);
      setTransformStatus(data);
      
      if (data.status === 'pending' || data.status === 'processing') {
        setTimeout(() => pollStatus(tfId), 2000);
      } else {
        setIsTransforming(false);
      }
    } catch (error) {
      console.error('Polling failed', error);
      setIsTransforming(false);
    }
  };

  if (isLoading) return <Loader fullPage size="lg" />;
  if (!image) return null;

  return (
    <div className="container mx-auto px-4 pb-12">
      <Button variant="ghost" onClick={() => navigate('/')} className="mb-6 -ml-4 text-gray-400 hover:text-white">
        <ArrowLeft className="w-4 h-4 mr-2" /> Back to Dashboard
      </Button>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <div className="lg:col-span-2 space-y-6">
          <Card>
            <CardHeader>
              <h2 className="text-xl font-semibold">Original Image Details</h2>
            </CardHeader>
            <CardContent>
              <div className="bg-bg-tertiary rounded-lg p-8 flex items-center justify-center min-h-[300px] border border-glass-border mb-6 relative overflow-hidden">
                <div className="absolute inset-0 opacity-10 bg-[radial-gradient(ellipse_at_center,_var(--tw-gradient-stops))] from-accent-primary via-transparent to-transparent"></div>
                <div className="text-center relative z-10">
                  <p className="text-xl font-mono text-gray-300">{image.original_filename}</p>
                  <p className="text-sm text-gray-500 mt-2">S3 Key: {image.s3_key}</p>
                </div>
              </div>
              
              <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div className="p-4 bg-white/5 rounded-lg border border-glass-border">
                  <p className="text-xs text-gray-500 mb-1">Format</p>
                  <p className="font-semibold uppercase">{image.format}</p>
                </div>
                <div className="p-4 bg-white/5 rounded-lg border border-glass-border">
                  <p className="text-xs text-gray-500 mb-1">Dimensions</p>
                  <p className="font-semibold">{image.width} x {image.height}</p>
                </div>
                <div className="p-4 bg-white/5 rounded-lg border border-glass-border">
                  <p className="text-xs text-gray-500 mb-1">Size</p>
                  <p className="font-semibold">{(image.size / 1024).toFixed(2)} KB</p>
                </div>
                <div className="p-4 bg-white/5 rounded-lg border border-glass-border">
                  <p className="text-xs text-gray-500 mb-1">Status</p>
                  <p className="font-semibold capitalize text-green-400">{image.status}</p>
                </div>
              </div>
            </CardContent>
          </Card>

          {transformStatus && (
            <Card className="animate-fade-in border-accent-primary/30">
              <CardHeader>
                <h3 className="text-lg font-semibold flex items-center">
                  Transformation Result
                  {transformStatus.status === 'completed' && <CheckCircle2 className="w-5 h-5 ml-2 text-green-500" />}
                  {transformStatus.status === 'failed' && <AlertCircle className="w-5 h-5 ml-2 text-red-500" />}
                  {(transformStatus.status === 'pending' || transformStatus.status === 'processing') && (
                    <RefreshCw className="w-5 h-5 ml-2 text-blue-500 animate-spin" />
                  )}
                </h3>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <p className="text-sm text-gray-500">Status</p>
                    <p className={`font-medium capitalize ${
                      transformStatus.status === 'completed' ? 'text-green-400' :
                      transformStatus.status === 'failed' ? 'text-red-400' : 'text-blue-400'
                    }`}>{transformStatus.status}</p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-500">Output S3 Key</p>
                    <p className="font-mono text-sm break-all">{transformStatus.s3_key || 'Processing...'}</p>
                  </div>
                </div>
              </CardContent>
            </Card>
          )}
        </div>

        <div>
          <Card className="sticky top-24">
            <CardHeader className="bg-accent-primary/10 border-b border-accent-primary/20">
              <h2 className="text-xl font-semibold flex items-center text-accent-primary">
                <Wand2 className="w-5 h-5 mr-2" /> Apply Transformations
              </h2>
            </CardHeader>
            <CardContent className="pt-6">
              <form onSubmit={handleTransform}>
                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-300 mb-1">Format</label>
                    <select 
                      value={format} 
                      onChange={(e) => setFormat(e.target.value)}
                      className="w-full bg-glass border border-glass-border rounded-lg px-4 py-2.5 text-white focus:outline-none focus:ring-2 focus:ring-accent-primary"
                    >
                      <option value="jpeg">JPEG</option>
                      <option value="png">PNG</option>
                      <option value="webp">WEBP</option>
                    </select>
                  </div>
                  
                  <Input 
                    label="Quality (1-100)" 
                    type="number" 
                    min="1" max="100" 
                    value={quality} 
                    onChange={(e) => setQuality(e.target.value)} 
                  />

                  <div className="grid grid-cols-2 gap-4">
                    <Input 
                      label="Resize Width" 
                      type="number" 
                      placeholder="e.g. 800" 
                      value={resizeW} 
                      onChange={(e) => setResizeW(e.target.value)} 
                    />
                    <Input 
                      label="Resize Height" 
                      type="number" 
                      placeholder="e.g. 600" 
                      value={resizeH} 
                      onChange={(e) => setResizeH(e.target.value)} 
                    />
                  </div>

                  <div className="flex items-center space-x-3 p-3 bg-white/5 rounded-lg border border-glass-border">
                    <input 
                      type="checkbox" 
                      id="grayscale" 
                      checked={grayscale} 
                      onChange={(e) => setGrayscale(e.target.checked)}
                      className="w-4 h-4 rounded border-gray-600 text-accent-primary focus:ring-accent-primary bg-bg-tertiary"
                    />
                    <label htmlFor="grayscale" className="text-sm font-medium text-gray-300 select-none">
                      Apply Grayscale Filter
                    </label>
                  </div>

                  <div className="flex items-center space-x-3 p-3 bg-white/5 rounded-lg border border-glass-border">
                    <input 
                      type="checkbox" 
                      id="flip" 
                      checked={flip} 
                      onChange={(e) => setFlip(e.target.checked)}
                      className="w-4 h-4 rounded border-gray-600 text-accent-primary focus:ring-accent-primary bg-bg-tertiary"
                    />
                    <label htmlFor="flip" className="text-sm font-medium text-gray-300 select-none">
                      Flip Vertically
                    </label>
                  </div>
                </div>

                <Button 
                  type="submit" 
                  className="w-full mt-8" 
                  isLoading={isTransforming}
                >
                  {isTransforming ? 'Processing...' : 'Start Transformation'}
                </Button>
              </form>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}

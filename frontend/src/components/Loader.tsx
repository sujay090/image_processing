import { Loader2 } from 'lucide-react';

interface LoaderProps {
  size?: 'sm' | 'md' | 'lg' | 'xl';
  className?: string;
  fullPage?: boolean;
}

export function Loader({ size = 'md', className = '', fullPage = false }: LoaderProps) {
  const sizes = {
    sm: 'w-4 h-4',
    md: 'w-8 h-8',
    lg: 'w-12 h-12',
    xl: 'w-16 h-16'
  };

  const loaderContent = (
    <Loader2 className={`${sizes[size]} animate-spin text-accent-primary ${className}`} />
  );

  if (fullPage) {
    return (
      <div className="fixed inset-0 flex items-center justify-center bg-primary/80 backdrop-blur-sm z-50">
        {loaderContent}
      </div>
    );
  }

  return (
    <div className="flex justify-center items-center p-4">
      {loaderContent}
    </div>
  );
}

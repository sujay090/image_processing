import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../lib/AuthContext';
import { Image as ImageIcon, LogOut } from 'lucide-react';
import { Button } from './Button';

export function Navbar() {
  const { isAuthenticated, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <nav className="sticky top-0 z-40 w-full glass-panel rounded-none border-t-0 border-l-0 border-r-0 mb-8">
      <div className="container mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          <div className="flex items-center">
            <Link to="/" className="flex items-center gap-2 group">
              <div className="p-2 bg-accent-primary/20 rounded-lg group-hover:bg-accent-primary/30 transition-colors">
                <ImageIcon className="h-6 w-6 text-accent-primary" />
              </div>
              <span className="font-bold text-xl tracking-tight text-white group-hover:text-transparent group-hover:bg-clip-text group-hover:bg-gradient-to-r group-hover:from-accent-primary group-hover:to-accent-secondary transition-all">
                VeoImages
              </span>
            </Link>
          </div>
          
          <div className="flex items-center gap-4">
            {isAuthenticated ? (
              <>
                <Link to="/" className="text-sm font-medium text-gray-300 hover:text-white transition-colors">
                  Dashboard
                </Link>
                <div className="h-4 w-px bg-glass-border"></div>
                <Button variant="ghost" size="sm" onClick={handleLogout} className="text-gray-300 hover:text-red-400">
                  <LogOut className="h-4 w-4 mr-2" />
                  Logout
                </Button>
              </>
            ) : (
              <>
                <Link to="/login">
                  <Button variant="ghost" size="sm">Login</Button>
                </Link>
                <Link to="/register">
                  <Button variant="primary" size="sm">Sign Up</Button>
                </Link>
              </>
            )}
          </div>
        </div>
      </div>
    </nav>
  );
}

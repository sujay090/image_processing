import React, { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../lib/AuthContext';
import { apiFetch } from '../lib/api';
import { Button } from '../components/Button';
import { Input } from '../components/Input';
import { Card, CardContent, CardHeader, CardFooter } from '../components/Card';
import { Image as ImageIcon } from 'lucide-react';

export function AuthPage() {
  const location = useLocation();
  const isInitialLogin = location.pathname === '/login';
  const [isLogin, setIsLogin] = useState(isInitialLogin);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  
  const { login } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError('');

    try {
      const endpoint = isLogin ? '/auth/login' : '/auth/register';
      const data = await apiFetch(endpoint, {
        method: 'POST',
        data: { email, password },
      });

      if (data.token) {
        login(data.token);
        navigate('/');
      } else {
        // Fallback if structure is different
        if (!isLogin) {
          // If registration doesn't return a token, maybe login is needed
          setIsLogin(true);
          setError('Registration successful. Please log in.');
        }
      }
    } catch (err: any) {
      setError(err?.error || err?.message || 'Authentication failed. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  const toggleMode = () => {
    setIsLogin(!isLogin);
    setError('');
    // Optionally push state to history so back button works
    navigate(isLogin ? '/register' : '/login', { replace: true });
  };

  return (
    <div className="flex flex-col items-center justify-center min-h-[80vh] px-4">
      <div className="mb-8 text-center animate-fade-in">
        <div className="inline-flex items-center justify-center p-3 bg-accent-primary/20 rounded-2xl mb-4 shadow-lg shadow-accent-primary/10">
          <ImageIcon className="w-10 h-10 text-accent-primary" />
        </div>
        <h1 className="text-3xl font-bold mb-2">
          Welcome to <span className="text-gradient">VeoImages</span>
        </h1>
        <p className="text-gray-400">
          Advanced asynchronous image processing
        </p>
      </div>

      <Card className="w-full max-w-md animate-fade-in shadow-[0_0_40px_rgba(99,102,241,0.1)]">
        <CardHeader>
          <h2 className="text-xl font-semibold text-center text-white">
            {isLogin ? 'Sign in to your account' : 'Create a new account'}
          </h2>
        </CardHeader>
        
        <CardContent>
          <form onSubmit={handleSubmit}>
            <Input
              label="Email Address"
              type="email"
              placeholder="you@example.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
            <Input
              label="Password"
              type="password"
              placeholder="••••••••"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
            
            {error && (
              <div className="mb-4 p-3 bg-red-500/10 border border-red-500/20 rounded-lg">
                <p className="text-sm text-red-400 text-center">{error}</p>
              </div>
            )}
            
            <Button 
              type="submit" 
              className="w-full mt-2"
              isLoading={isLoading}
            >
              {isLogin ? 'Sign In' : 'Create Account'}
            </Button>
          </form>
        </CardContent>
        
        <CardFooter className="justify-center">
          <p className="text-sm text-gray-400">
            {isLogin ? "Don't have an account? " : "Already have an account? "}
            <button 
              onClick={toggleMode}
              className="text-accent-primary hover:text-accent-secondary font-medium transition-colors focus:outline-none"
            >
              {isLogin ? 'Sign up' : 'Sign in'}
            </button>
          </p>
        </CardFooter>
      </Card>
    </div>
  );
}

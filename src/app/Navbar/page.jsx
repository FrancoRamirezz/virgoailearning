'use client'
import React, { useState } from 'react';
import { 
  BookOpen, 
  Play, 
  Phone,
  DollarSign,
  Info,
  Menu, 
  X, 
  User
} from 'lucide-react';

const Navbar = ({ currentView, setCurrentView, isAuthenticated, userRole }) => {
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  // Navigation items
  const navItems = [
    { id: 'home', label: 'Home', icon: BookOpen },
    { id: 'About', label: 'About', icon: Info },
    { id: 'Courses', label: 'Courses', icon: Play },
    { id: 'Payment', label: 'Pricing', icon: DollarSign },
    { id: 'Contact', label: 'Contact', icon: Phone },
  ];

  const handleNavClick = (itemId) => {
    console.log('Navbar: Attempting to navigate to:', itemId);
    console.log('Navbar: setCurrentView function exists:', typeof setCurrentView === 'function');
    
    // Check if setCurrentView is a function before calling it
    if (typeof setCurrentView === 'function') {
      setCurrentView(itemId);
    } else {
      console.error('Navbar: setCurrentView is not a function. Check if props are passed correctly.');
      console.log('Navbar: Received props:', { currentView, setCurrentView, isAuthenticated, userRole });
    }
    
    setIsMobileMenuOpen(false);
  };

  return (
    <nav className="bg-white/90 backdrop-blur-md shadow-lg border-b border-purple-200 sticky top-0 z-50">
      <div className="px-4 mx-auto max-w-7xl lg:px-8">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
          <div 
            className="flex items-center space-x-3 cursor-pointer"
            onClick={() => handleNavClick('home')}
          >
            <div className="w-10 h-10 bg-gradient-to-r from-purple-600 to-blue-600 rounded-xl flex items-center justify-center">
              <BookOpen className="w-6 h-6 text-white" />
            </div>
            <span className="text-xl font-bold text-gray-900">CitizenshipPrep</span>
          </div>
          
          {/* Desktop Navigation */}
          <div className="hidden md:flex items-center space-x-8">
            {navItems.map((item) => (
              <button
                key={item.id}
                onClick={() => handleNavClick(item.id)}
                className={`flex items-center space-x-2 px-3 py-2 rounded-lg transition-colors duration-200 ${
                  currentView === item.id 
                    ? 'text-purple-600 bg-purple-50 shadow-sm' 
                    : 'text-gray-600 hover:text-gray-900 hover:bg-gray-50'
                }`}
              >
                <item.icon className="w-4 h-4" />
                <span className="font-medium">{item.label}</span>
              </button>
            ))}
          </div>
          
          {/* Auth Buttons */}
          <div className="hidden md:flex items-center space-x-4">
            {!isAuthenticated ? (
              <>
                <button className="text-gray-600 hover:text-gray-900 font-medium transition-colors duration-200 px-3 py-2 rounded-lg hover:bg-gray-50">
                  Log In
                </button>
                <button className="bg-gradient-to-r from-purple-600 to-blue-600 text-white px-6 py-2 rounded-lg font-medium hover:from-purple-700 hover:to-blue-700 transition-all duration-300 shadow-lg hover:shadow-xl transform hover:scale-105">
                  Sign Up
                </button>
              </>
            ) : (
              <div className="flex items-center space-x-3 px-3 py-2 rounded-lg bg-gray-50">
                <div className="w-8 h-8 bg-gradient-to-r from-purple-600 to-blue-600 rounded-full flex items-center justify-center">
                  <User className="w-5 h-5 text-white" />
                </div>
                <span className="text-gray-700 font-medium capitalize">
                  {userRole || 'Student'}
                </span>
              </div>
            )}
          </div>
          
          {/* Mobile Menu Button */}
          <button
            className="md:hidden p-2 rounded-lg hover:bg-gray-100 transition-colors duration-200"
            onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
            aria-label="Toggle mobile menu"
          >
            {isMobileMenuOpen ? (
              <X className="w-6 h-6 text-gray-600" />
            ) : (
              <Menu className="w-6 h-6 text-gray-600" />
            )}
          </button>
        </div>
        
        {/* Mobile Menu */}
        {isMobileMenuOpen && (
          <div className="md:hidden py-4 border-t border-purple-200 bg-white/95 backdrop-blur-sm">
            <div className="flex flex-col space-y-1">
              {navItems.map((item) => (
                <button
                  key={item.id}
                  onClick={() => handleNavClick(item.id)}
                  className={`flex items-center space-x-3 px-4 py-3 rounded-lg text-left transition-colors duration-200 mx-2 ${
                    currentView === item.id 
                      ? 'text-purple-600 bg-purple-50 shadow-sm' 
                      : 'text-gray-600 hover:text-gray-900 hover:bg-gray-50'
                  }`}
                >
                  <item.icon className="w-5 h-5" />
                  <span className="font-medium">{item.label}</span>
                </button>
              ))}
              
              <div className="pt-4 mt-4 border-t border-purple-200 space-y-2 mx-2">
                {!isAuthenticated ? (
                  <>
                    <button className="w-full text-left px-4 py-3 text-gray-600 hover:text-gray-900 hover:bg-gray-50 rounded-lg transition-colors duration-200 font-medium">
                      Log In
                    </button>
                    <button className="w-full text-left px-4 py-3 bg-gradient-to-r from-purple-600 to-blue-600 text-white rounded-lg font-medium shadow-lg">
                      Sign Up
                    </button>
                  </>
                ) : (
                  <div className="px-4 py-3 flex items-center space-x-3 bg-gray-50 rounded-lg">
                    <div className="w-8 h-8 bg-gradient-to-r from-purple-600 to-blue-600 rounded-full flex items-center justify-center">
                      <User className="w-5 h-5 text-white" />
                    </div>
                    <span className="text-gray-700 font-medium capitalize">
                      {userRole || 'Student'}
                    </span>
                  </div>
                )}
              </div>
            </div>
          </div>
        )}
      </div>
    </nav>
  );
};

export default Navbar;
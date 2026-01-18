"use client";

import React, { useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  BookOpen,
  Play,
  Phone,
  DollarSign,
  Info,
  Menu,
  X,
  User,
} from "lucide-react";

interface NavbarProps {
  isAuthenticated?: boolean;
  userRole?: string;
}

interface NavItem {
  id: string;
  label: string;
  href: string;
  icon: React.ComponentType<React.SVGProps<SVGSVGElement>>;
}

const Navbar = ({ isAuthenticated = true, userRole }: NavbarProps) => {
  
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const pathname = usePathname();

  const navItems: NavItem[] = [
    { id: "home", label: "Dashboard", href: "/dashboard", icon: BookOpen },
    { id: "about", label: "About", href: "/about", icon: Info },
    { id: "courses", label: "Courses", href: "/courses", icon: Play },
    { id: "payment", label: "Pricing", href: "/payment", icon: DollarSign },
    { id: "contact", label: "Contact", href: "/contact", icon: Phone },
  ];

  const handleToggleMobileMenu = () => {
    setIsMobileMenuOpen((prev) => !prev);
  };

  const closeMobileMenu = () => {
    setIsMobileMenuOpen(false);
  };

  const isActive = (href: string) => {
    // Basic active check, can adjust if you use nested routes
    return pathname === href;
  };

  return (
    <nav className="bg-white/90 backdrop-blur-md shadow-lg border-b border-red-200 sticky top-0 z-50">
      <div className="px-4 mx-auto max-w-7xl ">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
          <Link
            href={navItems[0].href}
            className="flex items-center space-x-3 cursor-pointer"
            onClick={closeMobileMenu}
          >
            <div className="w-10 h-10 bg-gradient-to-r from-red-600 to-blue-600 rounded-xl flex items-center justify-center">
              <BookOpen className="w-6 h-6 text-white" />
            </div>
            <span className="text-xl font-bold text-gray-900">
              VirgoLearning
            </span>
          </Link>

          {/* Desktop Navigation */}
          <div className="hidden md:flex items-center ">
            {navItems.map((item) => (
              <Link
                key={item.id}
                href={item.href}
                className={`flex items-center space-x-2 px-3 py-2 rounded-lg transition-colors duration-200 ${
                  isActive(item.href)
                    ? "text-red-600 bg-red-50 shadow-sm"
                    : "text-gray-600 hover:text-gray-900 hover:bg-gray-50"
                }`}
              >
                <item.icon className="w-4 h-4" />
                <span className="font-medium">{item.label}</span>
              </Link>
            ))}
          </div>

          {/* Auth Buttons */}
          <div className="hidden md:flex items-center space-x-4">
            {!isAuthenticated ? (
              <>
                <Link
                  href="/login"
                  className="text-gray-600 hover:text-gray-900 font-medium transition-colors duration-200 px-3 py-2 rounded-lg hover:bg-gray-50"
                >
                  Log In
                </Link>
                <Link
                  href="/signup"
                  className="bg-gradient-to-r from-red-600 to-blue-600 text-white px-6 py-2 rounded-lg font-medium hover:from-red-700 hover:to-blue-700 transition-all duration-300 shadow-lg hover:shadow-xl transform hover:scale-105"
                >
                  Sign Up
                </Link>
              </>
            ) : (
              <div className="flex items-center space-x-3 px-3 py-2 rounded-lg bg-gray-50">
                <div className="w-8 h-8 bg-gradient-to-r from-red-600 to-blue-600 rounded-full flex items-center justify-center">
                  <User className="w-5 h-5 text-white" />
                </div>
                <span className="text-gray-700 font-medium capitalize">
                  {userRole || "Student"}
                </span>
              </div>
            )}
          </div>

          {/* Mobile Menu Button */}
          <button
            className="md:hidden p-2 rounded-lg hover:bg-gray-100 transition-colors duration-200"
            onClick={handleToggleMobileMenu}
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
          <div className="md:hidden py-4 border-t border-red-200 bg-white/95 backdrop-blur-sm">
            <div className="flex flex-col space-y-1">
              {navItems.map((item) => (
                <Link
                  key={item.id}
                  href={item.href}
                  onClick={closeMobileMenu}
                  className={`flex items-center space-x-3 px-4 py-3 rounded-lg text-left transition-colors duration-200 mx-2 ${
                    isActive(item.href)
                      ? "text-red-600 bg-red-50 shadow-sm"
                      : "text-gray-600 hover:text-gray-900 hover:bg-gray-50"
                  }`}
                >
                  <item.icon className="w-5 h-5" />
                  <span className="font-medium">{item.label}</span>
                </Link>
              ))}

              <div className="pt-4 mt-4 border-t border-red-200 space-y-2 mx-2">
                {!isAuthenticated ? (
                  <>
                    <Link
                      href="/login"
                      onClick={closeMobileMenu}
                      className="w-full text-left px-4 py-3 text-gray-600 hover:text-gray-900 hover:bg-gray-50 rounded-lg transition-colors duration-200 font-medium block"
                    >
                      Log In
                    </Link>
                    <Link
                      href="/signup"
                      onClick={closeMobileMenu}
                      className="w-full text-left px-4 py-3 bg-gradient-to-r from-red-600 to-blue-600 text-white rounded-lg font-medium shadow-lg block"
                    >
                      Sign Up
                    </Link>
                  </>
                ) : (
                  <div className="px-4 py-3 flex items-center space-x-3 bg-gray-50 rounded-lg">
                    <div className="w-8 h-8 bg-gradient-to-r from-red-600 to-blue-600 rounded-full flex items-center justify-center">
                      <User className="w-5 h-5 text-white" />
                    </div>
                    <span className="text-gray-700 font-medium capitalize">
                      {userRole || "Student"}
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

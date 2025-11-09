"use client"
import React, { useState } from 'react';
import { 
  BookOpen, 
  Calendar, 
  BarChart3, 
  Home, 
  FileText, 
  MessageCircle, 
  Settings,
  ChevronLeft,
  ChevronRight,
  User,
  Trophy,
  Clock
} from 'lucide-react';

const Sidebar = ({ activeSection, setActiveSection, isCollapsed, setIsCollapsed }) => {
  const [hoveredItem, setHoveredItem] = useState(null);

  const menuItems = [
    { id: 'dashboard', label: 'Dashboard', icon: Home, color: 'text-blue-600' },
    { id: 'courses', label: 'My Courses', icon: BookOpen, color: 'text-purple-600' },
    { id: 'calendar', label: 'Calendar', icon: Calendar, color: 'text-green-600' },
    { id: 'grades', label: 'Grades', icon: BarChart3, color: 'text-red-600' },
    { id: 'assignments', label: 'Assignments', icon: FileText, color: 'text-orange-600' },
    { id: 'discussions', label: 'Discussions', icon: MessageCircle, color: 'text-indigo-600' },
    { id: 'profile', label: 'Profile', icon: User, color: 'text-gray-600' },
    { id: 'settings', label: 'Settings', icon: Settings, color: 'text-gray-500' }
  ];

  const quickStats = [
    { label: 'Courses', value: '3', icon: BookOpen, color: 'bg-purple-100 text-purple-600' },
    { label: 'Hours Studied', value: '24', icon: Clock, color: 'bg-blue-100 text-blue-600' },
    { label: 'Achievements', value: '8', icon: Trophy, color: 'bg-yellow-100 text-yellow-600' }
  ];

  return (
    <div className={`${isCollapsed ? 'w-20' : 'w-72'} bg-white/90 backdrop-blur-sm border-r border-purple-200 h-screen sticky top-0 transition-all duration-300 flex flex-col shadow-lg`}>
      {/* Header */}
      <div className="p-6 border-b border-purple-100">
        <div className="flex items-center justify-between">
          {!isCollapsed && (
            <div>
              <h2 className="text-xl font-bold text-gray-900">Course Portal</h2>
              <p className="text-sm text-gray-600">Welcome back, Student!</p>
            </div>
          )}
          
          <button
            onClick={() => setIsCollapsed(!isCollapsed)}
            className="p-2 rounded-lg hover:bg-purple-100 transition-colors duration-200"
            aria-label="Toggle sidebar"
          >
            {isCollapsed ? (
              <ChevronRight className="w-5 h-5 text-gray-600" />
            ) : (
              <ChevronLeft className="w-5 h-5 text-gray-600" />
            )}
          </button>
        </div>
      </div>

      {/* Navigation Menu */}
      <div className="flex-1 py-6 overflow-y-auto">
        <nav className="space-y-2 px-4">
          {menuItems.map((item) => (
            <button
              key={item.id}
              onClick={() => setActiveSection(item.id)}
              onMouseEnter={() => setHoveredItem(item.id)}
              onMouseLeave={() => setHoveredItem(null)}
              className={`w-full flex items-center ${isCollapsed ? 'justify-center px-3' : 'px-4'} py-3 rounded-xl transition-all duration-200 group relative ${
                activeSection === item.id
                  ? 'bg-gradient-to-r from-purple-100 to-blue-100 text-purple-700 shadow-sm'
                  : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
              }`}
            >
              <item.icon className={`w-5 h-5 ${activeSection === item.id ? item.color : ''} group-hover:scale-110 transition-transform duration-200`} />
              
              {!isCollapsed && (
                <span className="ml-3 font-medium">{item.label}</span>
              )}

              {/* Tooltip for collapsed state */}
              {isCollapsed && hoveredItem === item.id && (
                <div className="absolute left-full ml-2 px-3 py-2 bg-gray-900 text-white text-sm rounded-lg whitespace-nowrap z-50">
                  {item.label}
                  <div className="absolute top-1/2 -left-1 transform -translate-y-1/2 w-2 h-2 bg-gray-900 rotate-45"></div>
                </div>
              )}
            </button>
          ))}
        </nav>

        {/* Quick Stats - Only show when expanded */}
        {!isCollapsed && (
          <div className="px-4 mt-8">
            <h3 className="text-sm font-semibold text-gray-500 uppercase tracking-wide mb-4">
              Quick Stats
            </h3>
            <div className="space-y-3">
              {quickStats.map((stat, index) => (
                <div key={index} className="flex items-center p-3 bg-gray-50 rounded-xl hover:bg-gray-100 transition-colors duration-200">
                  <div className={`w-8 h-8 rounded-lg ${stat.color} flex items-center justify-center`}>
                    <stat.icon className="w-4 h-4" />
                  </div>
                  <div className="ml-3">
                    <div className="text-lg font-bold text-gray-900">{stat.value}</div>
                    <div className="text-xs text-gray-500">{stat.label}</div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>

      {/* User Profile - Bottom */}
      <div className="p-4 border-t border-purple-100">
        {isCollapsed ? (
          <div className="flex justify-center">
            <div className="w-8 h-8 bg-gradient-to-r from-purple-600 to-blue-600 rounded-full flex items-center justify-center">
              <User className="w-4 h-4 text-white" />
            </div>
          </div>
        ) : (
          <div className="flex items-center p-3 rounded-xl hover:bg-gray-50 transition-colors duration-200 cursor-pointer">
            <div className="w-10 h-10 bg-gradient-to-r from-purple-600 to-blue-600 rounded-full flex items-center justify-center">
              <User className="w-5 h-5 text-white" />
            </div>
            <div className="ml-3">
              <div className="font-medium text-gray-900">John Doe</div>
              <div className="text-sm text-gray-500">Premium Student</div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default Sidebar;
"use client"
import React, { useState, useEffect } from 'react';
import Sidebar from './Sidebar';
import Calendar from './Calander';
import Grades from './grades';
import { BookOpen, TrendingUp, Clock, Star, Play, Trophy, Target, Zap, Award, ChevronRight, Bell, Search, Filter, Flame } from 'lucide-react';

const CoursePortal = () => {
  const [activeSection, setActiveSection] = useState('dashboard');
  const [isCollapsed, setIsCollapsed] = useState(false);
  const [selectedCourse, setSelectedCourse] = useState(null);
  const [hoveredCourse, setHoveredCourse] = useState(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [filterActive, setFilterActive] = useState('all');
  const [showNotification, setShowNotification] = useState(true);

  const dashboardData = {
    user: {
      name: 'Franco The Man',
      role: 'Learner',
      joinedDate: 'September 2024',
      streak: 7,
      totalXP: 2450,
      level: 'Intermediate'
    },
    enrolledCourses: [
      {
        id: 'english-foundations',
        title: 'English Language Foundations',
        university: 'American English Institute',
        instructor: 'Dr. Sarah Johnson',
        progress: 75,
        rating: 4.8,
        ratingCount: 12456,
        weeks: 8,
        hoursPerWeek: '4-6 hours',
        nextDeadline: 'Assignment due Jan 20',
        image: 'https://images.unsplash.com/photo-1434030216411-0b793f4b4173?w=300&h=200&fit=crop',
        level: 'Beginner',
        category: 'Language Learning',
        recentActivity: '2 hours ago',
        newContent: true
      },
      {
        id: 'business-english',
        title: 'Business English Communication',
        university: 'Professional Language Center',
        instructor: 'Prof. Michael Chen',
        progress: 45,
        rating: 4.7,
        ratingCount: 8934,
        weeks: 6,
        hoursPerWeek: '3-4 hours',
        nextDeadline: 'Quiz due Jan 18',
        image: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=300&h=200&fit=crop',
        level: 'Intermediate',
        category: 'Business',
        recentActivity: '1 day ago',
        newContent: false
      },
      {
        id: 'conversation-skills',
        title: 'Everyday English Conversation',
        university: 'Global Language Academy',
        instructor: 'Emily Rodriguez',
        progress: 20,
        rating: 4.9,
        ratingCount: 15672,
        weeks: 4,
        hoursPerWeek: '2-3 hours',
        nextDeadline: 'Discussion post due Jan 22',
        image: 'https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?w=300&h=200&fit=crop',
        level: 'Beginner',
        category: 'Speaking',
        recentActivity: '3 days ago',
        newContent: true
      }
    ],
    weeklyGoal: {
      target: 10,
      completed: 7,
      percentage: 70
    }
  };

  const DashboardView = () => (
    <div className="max-w-7xl mx-auto space-y-6">
      {/* Enhanced Header with Gamification */}
      <div className="bg-gradient-to-r from-purple-600 to-blue-600 rounded-2xl p-6 text-white shadow-lg">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold mb-2">
              Welcome back, {dashboardData.user.name}! 👋
            </h1>
            <p className="text-purple-100">You're on fire! Keep the momentum going.</p>
          </div>
          
          <div className="flex items-center gap-6">
            {/* Streak Counter */}
            <div className="text-center bg-white/10 backdrop-blur-sm rounded-xl px-6 py-3">
              <div className="flex items-center gap-2 mb-1">
                <Flame className="w-5 h-5 text-orange-400" />
                <span className="text-2xl font-bold">{dashboardData.user.streak}</span>
              </div>
              <div className="text-xs text-purple-200">Day Streak</div>
            </div>
            
            {/* XP Counter */}
            <div className="text-center bg-white/10 backdrop-blur-sm rounded-xl px-6 py-3">
              <div className="flex items-center gap-2 mb-1">
                <Zap className="w-5 h-5 text-yellow-400" />
                <span className="text-2xl font-bold">{dashboardData.user.totalXP}</span>
              </div>
              <div className="text-xs text-purple-200">Total XP</div>
            </div>
          </div>
        </div>

        {/* Weekly Goal Progress */}
        <div className="mt-6 bg-white/10 backdrop-blur-sm rounded-xl p-4">
          <div className="flex items-center justify-between mb-2">
            <span className="font-medium">Weekly Learning Goal</span>
            <span className="text-sm">{dashboardData.weeklyGoal.completed}/{dashboardData.weeklyGoal.target} hours</span>
          </div>
          <div className="w-full bg-white/20 rounded-full h-3">
            <div 
              className="bg-gradient-to-r from-green-400 to-emerald-500 h-3 rounded-full transition-all duration-700 ease-out"
              style={{ width: `${dashboardData.weeklyGoal.percentage}%` }}
            ></div>
          </div>
        </div>
      </div>

      {/* Quick Actions Bar */}
      <div className="flex items-center justify-between bg-white rounded-xl p-4 shadow-md">
        <div className="flex items-center gap-3 flex-1">
          <Search className="w-5 h-5 text-gray-400" />
          <input
            type="text"
            placeholder="Search courses, lessons, or topics..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="flex-1 outline-none text-gray-700"
          />
        </div>
        
        <div className="flex items-center gap-3 border-l pl-4">
          <button
            onClick={() => setFilterActive(filterActive === 'all' ? 'active' : 'all')}
            className={`flex items-center gap-2 px-4 py-2 rounded-lg transition-all duration-200 ${
              filterActive === 'active' 
                ? 'bg-purple-100 text-purple-700' 
                : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
            }`}
          >
            <Filter className="w-4 h-4" />
            <span className="text-sm font-medium">Filter</span>
          </button>
          
          <button className="relative p-2 hover:bg-gray-100 rounded-lg transition-colors">
            <Bell className="w-5 h-5 text-gray-600" />
            {showNotification && (
              <div className="absolute top-1 right-1 w-2 h-2 bg-red-500 rounded-full"></div>
            )}
          </button>
        </div>
      </div>

      {/* Continue Learning - Enhanced */}
      <div className="bg-gradient-to-r from-blue-50 to-purple-50 rounded-2xl p-6 border-2 border-blue-200 shadow-md">
        <div className="flex items-center gap-2 mb-4">
          <Play className="w-5 h-5 text-blue-600" />
          <h2 className="text-xl font-bold text-gray-900">Continue Your Journey</h2>
        </div>
        
        <div className="bg-white rounded-xl p-6 hover:shadow-lg transition-all duration-300 cursor-pointer transform hover:scale-[1.02]">
          <div className="flex items-center justify-between">
            <div className="flex-1">
              <div className="flex items-center gap-2 mb-2">
                <span className="px-2 py-1 bg-green-100 text-green-700 text-xs font-semibold rounded">In Progress</span>
                <span className="text-sm text-gray-500">15 minutes left</span>
              </div>
              <h3 className="text-lg font-semibold text-gray-900 mb-1">
                Present Perfect Tense - Advanced Usage
              </h3>
              <p className="text-gray-600 mb-3">Grammar Fundamentals · English Language Foundations</p>
              
              <div className="w-full bg-gray-200 rounded-full h-2 mb-2">
                <div className="bg-gradient-to-r from-blue-500 to-purple-500 h-2 rounded-full" style={{ width: '75%' }}></div>
              </div>
              <p className="text-sm text-gray-500">75% complete</p>
            </div>
            
            <button className="ml-6 px-8 py-3 bg-gradient-to-r from-blue-600 to-purple-600 text-white rounded-xl font-semibold hover:from-blue-700 hover:to-purple-700 transition-all duration-200 transform hover:scale-105 shadow-lg flex items-center gap-2">
              Continue
              <ChevronRight className="w-5 h-5" />
            </button>
          </div>
        </div>
      </div>

      {/* My Courses - Interactive Grid */}
      <div>
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-2xl font-bold text-gray-900">My Courses</h2>
          <button 
            onClick={() => setActiveSection('courses')}
            className="text-purple-600 hover:text-purple-700 font-semibold flex items-center gap-1 group"
          >
            View all
            <ChevronRight className="w-4 h-4 group-hover:translate-x-1 transition-transform" />
          </button>
        </div>
        
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {dashboardData.enrolledCourses.map((course) => (
            <div 
              key={course.id} 
              className="group bg-white rounded-2xl overflow-hidden shadow-md hover:shadow-2xl transition-all duration-300 cursor-pointer transform hover:scale-[1.03]"
              onClick={() => {
                setSelectedCourse(course.id);
                setActiveSection('course');
              }}
              onMouseEnter={() => setHoveredCourse(course.id)}
              onMouseLeave={() => setHoveredCourse(null)}
            >
              <div className="relative overflow-hidden">
                <img 
                  src={course.image} 
                  alt={course.title}
                  className="w-full h-48 object-cover transform group-hover:scale-110 transition-transform duration-500"
                />
                
                {/* Overlay on Hover */}
                <div className={`absolute inset-0 bg-gradient-to-t from-black/60 to-transparent transition-opacity duration-300 ${
                  hoveredCourse === course.id ? 'opacity-100' : 'opacity-0'
                }`}>
                  <div className="absolute bottom-4 left-4 right-4">
                    <button className="w-full py-2 bg-white text-purple-600 rounded-lg font-semibold flex items-center justify-center gap-2">
                      <Play className="w-4 h-4" />
                      Resume Learning
                    </button>
                  </div>
                </div>
                
                {/* Badges */}
                <div className="absolute top-3 left-3 flex gap-2">
                  <span className="bg-white/90 backdrop-blur-sm px-3 py-1 rounded-full text-xs font-semibold text-gray-700">
                    {course.category}
                  </span>
                  {course.newContent && (
                    <span className="bg-gradient-to-r from-purple-500 to-pink-500 text-white px-3 py-1 rounded-full text-xs font-semibold animate-pulse">
                      New
                    </span>
                  )}
                </div>
              </div>
              
              <div className="p-5">
                <h3 className="font-bold text-gray-900 mb-2 line-clamp-2 text-lg group-hover:text-purple-600 transition-colors">
                  {course.title}
                </h3>
                <p className="text-sm text-gray-600 mb-3">{course.university}</p>
                
                <div className="flex items-center justify-between mb-4 text-sm">
                  <div className="flex items-center gap-1">
                    <Star className="w-4 h-4 text-yellow-400 fill-current" />
                    <span className="font-semibold">{course.rating}</span>
                    <span className="text-gray-500">({course.ratingCount.toLocaleString()})</span>
                  </div>
                  <span className="px-2 py-1 bg-purple-100 text-purple-700 rounded text-xs font-medium">
                    {course.level}
                  </span>
                </div>
                
                {/* Progress Bar */}
                <div className="mb-4">
                  <div className="flex justify-between text-sm mb-2">
                    <span className="text-gray-600 font-medium">{course.progress}% complete</span>
                    <span className="text-gray-500">{course.recentActivity}</span>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-2 overflow-hidden">
                    <div 
                      className="bg-gradient-to-r from-purple-500 to-blue-500 h-2 rounded-full transition-all duration-700 ease-out"
                      style={{ width: `${course.progress}%` }}
                    ></div>
                  </div>
                </div>
                
                <div className="flex items-center justify-between pt-3 border-t border-gray-100">
                  <span className="text-sm text-gray-600">
                    <Clock className="w-4 h-4 inline mr-1" />
                    {course.hoursPerWeek}
                  </span>
                  <span className="text-sm font-medium text-purple-600">
                    {course.nextDeadline}
                  </span>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Achievements Section */}
      <div className="bg-gradient-to-r from-yellow-50 to-orange-50 rounded-2xl p-6 border-2 border-yellow-200">
        <div className="flex items-center gap-2 mb-4">
          <Trophy className="w-6 h-6 text-yellow-600" />
          <h2 className="text-xl font-bold text-gray-900">Recent Achievements</h2>
        </div>
        
        <div className="grid grid-cols-3 gap-4">
          {[
            { title: '7-Day Streak', icon: Flame, color: 'from-orange-500 to-red-500' },
            { title: 'Fast Learner', icon: Zap, color: 'from-yellow-500 to-orange-500' },
            { title: 'Perfect Score', icon: Award, color: 'from-purple-500 to-pink-500' }
          ].map((achievement, index) => (
            <div key={index} className="bg-white rounded-xl p-4 text-center hover:shadow-lg transition-all duration-300 cursor-pointer transform hover:scale-105">
              <div className={`w-16 h-16 bg-gradient-to-r ${achievement.color} rounded-full flex items-center justify-center mx-auto mb-3`}>
                <achievement.icon className="w-8 h-8 text-white" />
              </div>
              <p className="font-semibold text-gray-900">{achievement.title}</p>
            </div>
          ))}
        </div>
      </div>
    </div>
  );

  const CoursesListView = () => (
    <div className="max-w-7xl mx-auto">
      <div className="bg-white rounded-xl shadow-md p-6 mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">My Learning</h1>
        <div className="flex items-center gap-6 mt-6">
          <button className="text-purple-600 border-b-2 border-purple-600 pb-3 font-semibold">
            All courses ({dashboardData.enrolledCourses.length})
          </button>
          <button className="text-gray-600 hover:text-gray-900 pb-3 font-medium">
            In Progress
          </button>
          <button className="text-gray-600 hover:text-gray-900 pb-3 font-medium">
            Completed
          </button>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {dashboardData.enrolledCourses.map((course) => (
          <div 
            key={course.id}
            className="bg-white rounded-2xl overflow-hidden shadow-md hover:shadow-xl transition-all duration-300 cursor-pointer transform hover:scale-105"
            onClick={() => {
              setSelectedCourse(course.id);
              setActiveSection('course');
            }}
          >
            <img 
              src={course.image} 
              alt={course.title}
              className="w-full h-36 object-cover"
            />
            
            <div className="p-4">
              <h3 className="font-semibold text-gray-900 mb-2 line-clamp-2">
                {course.title}
              </h3>
              
              <p className="text-xs text-gray-600 mb-3">{course.university}</p>
              
              <div className="mb-3">
                <div className="flex justify-between text-xs text-gray-600 mb-1">
                  <span className="font-medium">{course.progress}% complete</span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <div 
                    className="bg-gradient-to-r from-green-400 to-emerald-500 h-2 rounded-full transition-all duration-500"
                    style={{ width: `${course.progress}%` }}
                  ></div>
                </div>
              </div>
              
              <button className="w-full py-2 text-sm text-white bg-purple-600 rounded-lg hover:bg-purple-700 transition-colors font-medium">
                Continue
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );

  const renderMainContent = () => {
    switch (activeSection) {
      case 'dashboard':
        return <DashboardView />;
      case 'courses':
        return <CoursesListView />;
      case 'calendar':
        return <Calendar />;
      case 'grades':
        return <Grades />;
      default:
        return <DashboardView />;
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-purple-50 via-white to-blue-50 flex">
      <Sidebar
        activeSection={activeSection}
        setActiveSection={setActiveSection}
        isCollapsed={isCollapsed}
        setIsCollapsed={setIsCollapsed}
      />
      
      <main className="flex-1 p-8 overflow-y-auto">
        {renderMainContent()}
      </main>
    </div>
  );
};

export default CoursePortal;
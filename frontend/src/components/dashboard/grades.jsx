"use client"
import React, { useState } from 'react';
import { TrendingUp, TrendingDown, Award, Target, BarChart3, BookOpen, Clock, CheckCircle, AlertCircle } from 'lucide-react';

const Grades = () => {
  const [selectedCourse, setSelectedCourse] = useState('all');
  const [viewType, setViewType] = useState('overview'); // 'overview', 'detailed', 'progress'

  // Sample grades data
  const courses = [
    {
      id: 'citizenship-prep',
      name: 'Citizenship Test Preparation',
      instructor: 'Sarah Johnson',
      currentGrade: 92,
      totalPoints: 450,
      earnedPoints: 414,
      color: 'from-purple-600 to-blue-600',
      status: 'excellent'
    },
    {
      id: 'english-citizenship',
      name: 'English for Citizenship',
      instructor: 'Michael Chen',
      currentGrade: 88,
      totalPoints: 300,
      earnedPoints: 264,
      color: 'from-blue-600 to-indigo-600',
      status: 'good'
    },
    {
      id: 'interview-prep',
      name: 'Mock Interview Sessions',
      instructor: 'Emily Rodriguez',
      currentGrade: 85,
      totalPoints: 200,
      earnedPoints: 170,
      color: 'from-green-600 to-teal-600',
      status: 'good'
    }
  ];

  const assignments = [
    {
      id: 1,
      title: 'American History Quiz',
      course: 'Citizenship Test Preparation',
      type: 'quiz',
      score: 95,
      maxScore: 100,
      submittedDate: '2025-01-10',
      dueDate: '2025-01-12',
      status: 'graded',
      feedback: 'Excellent work! Strong understanding of key historical events.'
    },
    {
      id: 2,
      title: 'Government Structure Essay',
      course: 'Citizenship Test Preparation',
      type: 'essay',
      score: 88,
      maxScore: 100,
      submittedDate: '2025-01-08',
      dueDate: '2025-01-10',
      status: 'graded',
      feedback: 'Good analysis, but could use more specific examples.'
    },
    {
      id: 3,
      title: 'English Comprehension Test',
      course: 'English for Citizenship',
      type: 'test',
      score: 92,
      maxScore: 100,
      submittedDate: '2025-01-05',
      dueDate: '2025-01-07',
      status: 'graded',
      feedback: 'Great improvement in reading comprehension skills!'
    },
    {
      id: 4,
      title: 'Practice Interview Recording',
      course: 'Mock Interview Sessions',
      type: 'assignment',
      score: null,
      maxScore: 50,
      submittedDate: '2025-01-15',
      dueDate: '2025-01-16',
      status: 'submitted',
      feedback: null
    },
    {
      id: 5,
      title: 'Civics Knowledge Assessment',
      course: 'Citizenship Test Preparation',
      type: 'quiz',
      score: null,
      maxScore: 75,
      submittedDate: null,
      dueDate: '2025-01-20',
      status: 'pending',
      feedback: null
    }
  ];

  const getStatusColor = (status) => {
    switch (status) {
      case 'excellent': return 'text-green-600 bg-green-100';
      case 'good': return 'text-blue-600 bg-blue-100';
      case 'fair': return 'text-yellow-600 bg-yellow-100';
      case 'poor': return 'text-red-600 bg-red-100';
      default: return 'text-gray-600 bg-gray-100';
    }
  };

  const getGradeColor = (score) => {
    if (score >= 90) return 'text-green-600';
    if (score >= 80) return 'text-blue-600';
    if (score >= 70) return 'text-yellow-600';
    return 'text-red-600';
  };

  const getOverallGPA = () => {
    const totalWeighted = courses.reduce((sum, course) => sum + course.currentGrade, 0);
    return (totalWeighted / courses.length).toFixed(1);
  };

  const filteredAssignments = selectedCourse === 'all' 
    ? assignments 
    : assignments.filter(assignment => 
        courses.find(course => course.name === assignment.course)?.id === selectedCourse
      );

  const AssignmentStatusIcon = ({ status }) => {
    switch (status) {
      case 'graded':
        return <CheckCircle className="w-4 h-4 text-green-500" />;
      case 'submitted':
        return <Clock className="w-4 h-4 text-blue-500" />;
      case 'pending':
        return <AlertCircle className="w-4 h-4 text-orange-500" />;
      default:
        return <AlertCircle className="w-4 h-4 text-gray-500" />;
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="bg-white/90 backdrop-blur-sm rounded-2xl shadow-lg border border-purple-100 p-6">
        <div className="flex items-center justify-between mb-6">
          <div>
            <h2 className="text-2xl font-bold text-gray-900">Grades & Progress</h2>
            <p className="text-gray-600">Track your academic performance across all courses</p>
          </div>
          
          <div className="flex items-center space-x-4">
            <select
              value={selectedCourse}
              onChange={(e) => setSelectedCourse(e.target.value)}
              className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
            >
              <option value="all">All Courses</option>
              {courses.map(course => (
                <option key={course.id} value={course.id}>{course.name}</option>
              ))}
            </select>
            
            <div className="flex bg-gray-100 rounded-lg p-1">
              {['overview', 'detailed', 'progress'].map((type) => (
                <button
                  key={type}
                  onClick={() => setViewType(type)}
                  className={`px-4 py-2 rounded-md capitalize transition-all duration-200 ${
                    viewType === type
                      ? 'bg-white text-purple-600 shadow-sm'
                      : 'text-gray-600 hover:text-gray-900'
                  }`}
                >
                  {type}
                </button>
              ))}
            </div>
          </div>
        </div>

        {/* Overall Stats */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-6">
          <div className="text-center">
            <div className="w-16 h-16 bg-gradient-to-r from-purple-100 to-blue-100 rounded-xl flex items-center justify-center mx-auto mb-3">
              <Award className="w-8 h-8 text-purple-600" />
            </div>
            <div className="text-2xl font-bold text-gray-900">{getOverallGPA()}</div>
            <div className="text-sm text-gray-600">Overall GPA</div>
          </div>
          
          <div className="text-center">
            <div className="w-16 h-16 bg-gradient-to-r from-green-100 to-teal-100 rounded-xl flex items-center justify-center mx-auto mb-3">
              <TrendingUp className="w-8 h-8 text-green-600" />
            </div>
            <div className="text-2xl font-bold text-gray-900">95%</div>
            <div className="text-sm text-gray-600">Completion Rate</div>
          </div>
          
          <div className="text-center">
            <div className="w-16 h-16 bg-gradient-to-r from-blue-100 to-indigo-100 rounded-xl flex items-center justify-center mx-auto mb-3">
              <Target className="w-8 h-8 text-blue-600" />
            </div>
            <div className="text-2xl font-bold text-gray-900">8/10</div>
            <div className="text-sm text-gray-600">Goals Met</div>
          </div>
          
          <div className="text-center">
            <div className="w-16 h-16 bg-gradient-to-r from-orange-100 to-red-100 rounded-xl flex items-center justify-center mx-auto mb-3">
              <BarChart3 className="w-8 h-8 text-orange-600" />
            </div>
            <div className="text-2xl font-bold text-gray-900">24h</div>
            <div className="text-sm text-gray-600">Study Time</div>
          </div>
        </div>
      </div>

      {/* Course Grades */}
      {viewType === 'overview' && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {courses.map((course) => (
            <div key={course.id} className="bg-white/90 backdrop-blur-sm rounded-2xl shadow-lg border border-purple-100 p-6 hover:shadow-xl transition-all duration-300">
              <div className="flex items-center justify-between mb-4">
                <div className={`w-12 h-12 bg-gradient-to-r ${course.color} rounded-xl flex items-center justify-center`}>
                  <BookOpen className="w-6 h-6 text-white" />
                </div>
                <div className={`px-3 py-1 rounded-full text-xs font-semibold ${getStatusColor(course.status)}`}>
                  {course.status}
                </div>
              </div>
              
              <h3 className="text-lg font-semibold text-gray-900 mb-2">{course.name}</h3>
              <p className="text-sm text-gray-600 mb-4">Instructor: {course.instructor}</p>
              
              <div className="mb-4">
                <div className="flex items-center justify-between mb-2">
                  <span className="text-sm text-gray-600">Current Grade</span>
                  <span className={`text-xl font-bold ${getGradeColor(course.currentGrade)}`}>
                    {course.currentGrade}%
                  </span>
                </div>
                
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <div
                    className={`h-2 rounded-full bg-gradient-to-r ${course.color}`}
                    style={{ width: `${course.currentGrade}%` }}
                  ></div>
                </div>
              </div>
              
              <div className="text-sm text-gray-600">
                Points: {course.earnedPoints}/{course.totalPoints}
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Detailed View */}
      {viewType === 'detailed' && (
        <div className="bg-white/90 backdrop-blur-sm rounded-2xl shadow-lg border border-purple-100 p-6">
          <h3 className="text-xl font-semibold text-gray-900 mb-6">Assignment Details</h3>
          
          <div className="space-y-4">
            {filteredAssignments.map((assignment) => (
              <div key={assignment.id} className="border border-gray-200 rounded-xl p-4 hover:shadow-md transition-shadow duration-200">
                <div className="flex items-center justify-between mb-3">
                  <div className="flex items-center space-x-3">
                    <AssignmentStatusIcon status={assignment.status} />
                    <div>
                      <h4 className="font-semibold text-gray-900">{assignment.title}</h4>
                      <p className="text-sm text-gray-600">{assignment.course}</p>
                    </div>
                  </div>
                  
                  <div className="text-right">
                    {assignment.score !== null ? (
                      <div className={`text-lg font-bold ${getGradeColor((assignment.score / assignment.maxScore) * 100)}`}>
                        {assignment.score}/{assignment.maxScore}
                      </div>
                    ) : (
                      <div className="text-gray-500">
                        {assignment.status === 'pending' ? 'Not Submitted' : 'Pending Grade'}
                      </div>
                    )}
                    <div className="text-xs text-gray-500 capitalize">
                      {assignment.type}
                    </div>
                  </div>
                </div>
                
                <div className="flex items-center justify-between text-sm text-gray-600 mb-2">
                  <div>Due: {new Date(assignment.dueDate).toLocaleDateString()}</div>
                  {assignment.submittedDate && (
                    <div>Submitted: {new Date(assignment.submittedDate).toLocaleDateString()}</div>
                  )}
                </div>
                
                {assignment.feedback && (
                  <div className="mt-3 p-3 bg-blue-50 rounded-lg">
                    <div className="text-sm font-medium text-blue-800 mb-1">Instructor Feedback:</div>
                    <div className="text-sm text-blue-700">{assignment.feedback}</div>
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Progress View */}
      {viewType === 'progress' && (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Grade Trends */}
          <div className="bg-white/90 backdrop-blur-sm rounded-2xl shadow-lg border border-purple-100 p-6">
            <h3 className="text-xl font-semibold text-gray-900 mb-6">Grade Trends</h3>
            
            <div className="space-y-4">
              {courses.map((course) => (
                <div key={course.id} className="p-4 bg-gray-50 rounded-xl">
                  <div className="flex items-center justify-between mb-2">
                    <span className="font-medium text-gray-900">{course.name}</span>
                    <div className="flex items-center space-x-2">
                      <TrendingUp className="w-4 h-4 text-green-500" />
                      <span className={`font-bold ${getGradeColor(course.currentGrade)}`}>
                        {course.currentGrade}%
                      </span>
                    </div>
                  </div>
                  
                  <div className="w-full bg-gray-200 rounded-full h-2">
                    <div
                      className={`h-2 rounded-full bg-gradient-to-r ${course.color}`}
                      style={{ width: `${course.currentGrade}%` }}
                    ></div>
                  </div>
                </div>
              ))}
            </div>
          </div>
          
          {/* Recent Activity */}
          <div className="bg-white/90 backdrop-blur-sm rounded-2xl shadow-lg border border-purple-100 p-6">
            <h3 className="text-xl font-semibold text-gray-900 mb-6">Recent Activity</h3>
            
            <div className="space-y-4">
              {assignments.slice(0, 5).map((assignment) => (
                <div key={assignment.id} className="flex items-center space-x-3 p-3 rounded-lg hover:bg-gray-50 transition-colors duration-200">
                  <AssignmentStatusIcon status={assignment.status} />
                  <div className="flex-1">
                    <div className="font-medium text-gray-900">{assignment.title}</div>
                    <div className="text-sm text-gray-600">{assignment.course}</div>
                  </div>
                  <div className="text-right">
                    {assignment.score !== null ? (
                      <div className={`font-bold ${getGradeColor((assignment.score / assignment.maxScore) * 100)}`}>
                        {Math.round((assignment.score / assignment.maxScore) * 100)}%
                      </div>
                    ) : (
                      <div className="text-xs text-gray-500 capitalize">
                        {assignment.status}
                      </div>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Grades;
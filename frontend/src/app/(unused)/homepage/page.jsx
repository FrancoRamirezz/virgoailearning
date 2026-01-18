'use client'
import React, { useState } from 'react';
import { 
  BookOpen, 
  Users, 
  Award, 
  Star, 
  Play, 
  Clock, 
  CheckCircle, 
  Menu, 
  X, 
  User,
  Mail,
  Phone,
  MapPin,
  Send,
  DollarSign,
  Crown,
  Zap,
  Shield
} from 'lucide-react';

const LearningManagementSystem = () => {
  const [currentView, setCurrentView] = useState('home');
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const [userRole, setUserRole] = useState('student');
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  // Navigation items
  const navItems = [
    { id: 'home', label: 'Home', icon: BookOpen },
    { id: 'courses', label: 'Courses', icon: Play },
    { id: 'pricing', label: 'Pricing', icon: DollarSign },
    { id: 'contact', label: 'Contact', icon: Phone },
  ];

  // Sample course data
  const courses = [
    {
      id: 1,
      title: 'Citizenship Test Preparation',
      description: 'Complete preparation for the U.S. citizenship test with interactive lessons and practice tests.',
      instructor: 'Sarah Johnson',
      rating: 4.8,
      students: 2456,
      duration: '6 weeks',
      level: 'Beginner',
      price: 49.99,
      image: 'https://images.unsplash.com/photo-1521737604893-d14cc237f11d?w=400&h=300&fit=crop',
      modules: ['American History', 'Government Structure', 'Civics & Rights', 'Practice Tests']
    },
    {
      id: 2,
      title: 'Advanced English for Citizenship',
      description: 'Improve your English skills specifically for the citizenship interview and test.',
      instructor: 'Michael Chen',
      rating: 4.9,
      students: 1834,
      duration: '4 weeks',
      level: 'Intermediate',
      price: 39.99,
      image: 'https://images.unsplash.com/photo-1434030216411-0b793f4b4173?w=400&h=300&fit=crop',
      modules: ['Interview Preparation', 'Speaking Practice', 'Reading Comprehension', 'Writing Skills']
    },
    {
      id: 3,
      title: 'Mock Interview Sessions',
      description: 'Practice citizenship interviews with experienced instructors and AI-powered feedback.',
      instructor: 'Emily Rodriguez',
      rating: 4.7,
      students: 956,
      duration: '2 weeks',
      level: 'All Levels',
      price: 29.99,
      image: 'https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?w=400&h=300&fit=crop',
      modules: ['Interview Scenarios', 'Common Questions', 'Confidence Building', 'Feedback Sessions']
    }
  ];

  const pricingPlans = [
    {
      name: 'Basic',
      price: 29,
      period: 'month',
      features: [
        'Access to 5 courses',
        'Practice tests',
        'Email support',
        'Progress tracking',
        'Mobile app access'
      ],
      recommended: false,
      icon: Shield
    },
    {
      name: 'Pro',
      price: 49,
      period: 'month',
      features: [
        'Access to all courses',
        'Live tutoring sessions',
        'Priority support',
        'Advanced analytics',
        'Mock interviews',
        'Personalized study plan'
      ],
      recommended: true,
      icon: Crown
    },
    {
      name: 'Premium',
      price: 89,
      period: 'month',
      features: [
        'Everything in Pro',
        'One-on-one coaching',
        '24/7 support',
        'Custom learning path',
        'Guaranteed results',
        'Exam fee refund*'
      ],
      recommended: false,
      icon: Zap
    }
  ];

  const HomePage = () => (
    <div className="min-h-screen bg-gradient-to-b from-purple-100 to-blue-100">
      {/* Hero Section */}
      <section className="relative px-4 py-20 mx-auto max-w-7xl lg:px-8">
        <div className="grid items-center gap-12 lg:grid-cols-2">
          <div className="space-y-8">
            <div className="inline-flex items-center px-4 py-2 text-sm font-medium text-purple-700 bg-purple-200 rounded-full">
              <Star className="w-4 h-4 mr-2" />
              #1 Citizenship Test Prep Platform
            </div>
            
            <h1 className="text-5xl font-bold leading-tight text-gray-900 lg:text-6xl">
              Master Your
              <span className="text-transparent bg-gradient-to-r from-purple-600 to-blue-600 bg-clip-text"> U.S. Citizenship </span>
              Test
            </h1>
            
            <p className="text-xl text-gray-700 leading-relaxed">
              Join thousands of successful students who passed their citizenship test with our AI-powered learning platform. 
              Get personalized study plans, practice tests, and expert guidance.
            </p>
            
            <div className="flex flex-col gap-4 sm:flex-row">
              <button 
                onClick={() => setCurrentView('courses')}
                className="px-8 py-4 text-lg font-semibold text-white transition-all duration-300 bg-gradient-to-r from-purple-600 to-blue-600 rounded-xl hover:from-purple-700 hover:to-blue-700 transform hover:scale-105 shadow-lg hover:shadow-xl"
              >
                Start Learning Today
              </button>
              <button className="px-8 py-4 text-lg font-semibold text-gray-700 transition-all duration-300 bg-white border-2 border-purple-300 rounded-xl hover:border-purple-500 hover:text-purple-600 shadow-lg hover:shadow-xl">
                Try Free Demo
              </button>
            </div>
            
            {/* Stats */}
            <div className="flex gap-8 pt-8 border-t border-purple-200">
              <div className="text-center">
                <div className="text-3xl font-bold text-gray-900">98%</div>
                <div className="text-sm text-gray-600">Pass Rate</div>
              </div>
              <div className="text-center">
                <div className="text-3xl font-bold text-gray-900">15k+</div>
                <div className="text-sm text-gray-600">Students</div>
              </div>
              <div className="text-center">
                <div className="text-3xl font-bold text-gray-900">4.9★</div>
                <div className="text-sm text-gray-600">Rating</div>
              </div>
            </div>
          </div>
          
          <div className="relative">
            <div className="absolute inset-0 bg-gradient-to-r from-purple-400 to-blue-500 rounded-3xl rotate-3 opacity-30"></div>
            <img 
              src="https://images.unsplash.com/photo-1522202176988-66273c2fd55f?w=600&h=400&fit=crop"
              alt="Students learning"
              className="relative z-10 w-full h-96 object-cover rounded-3xl shadow-2xl"
            />
            <div className="absolute -bottom-6 -right-6 z-20 bg-white p-6 rounded-2xl shadow-xl border border-purple-100">
              <div className="flex items-center space-x-3">
                <div className="w-12 h-12 bg-gradient-to-r from-purple-100 to-blue-100 rounded-full flex items-center justify-center">
                  <CheckCircle className="w-6 h-6 text-green-600" />
                </div>
                <div>
                  <div className="font-semibold text-gray-900">Certificate Ready</div>
                  <div className="text-sm text-gray-600">in 6 weeks</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-20 bg-white/50 backdrop-blur-sm">
        <div className="px-4 mx-auto max-w-7xl lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-4xl font-bold text-gray-900 mb-4">
              Everything You Need to Succeed
            </h2>
            <p className="text-xl text-gray-700 max-w-3xl mx-auto">
              Our comprehensive platform combines expert instruction, AI-powered personalization, 
              and proven study methods to ensure your success.
            </p>
          </div>
          
          <div className="grid gap-8 md:grid-cols-3">
            {[
              {
                icon: BookOpen,
                title: 'Expert Curriculum',
                description: 'Comprehensive courses designed by citizenship test experts and immigration lawyers.'
              },
              {
                icon: Users,
                title: 'Live Tutoring',
                description: 'Get personalized help from certified instructors through live video sessions.'
              },
              {
                icon: Award,
                title: 'Practice Tests',
                description: 'Unlimited practice tests with detailed explanations and progress tracking.'
              }
            ].map((feature, index) => (
              <div key={index} className="p-8 bg-white/80 backdrop-blur-sm rounded-2xl hover:shadow-lg transition-all duration-300 group border border-purple-100 hover:border-purple-200">
                <div className="w-14 h-14 bg-gradient-to-r from-purple-600 to-blue-600 rounded-xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform duration-300">
                  <feature.icon className="w-7 h-7 text-white" />
                </div>
                <h3 className="text-xl font-semibold text-gray-900 mb-3">{feature.title}</h3>
                <p className="text-gray-700 leading-relaxed">{feature.description}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-20 bg-gradient-to-r from-purple-600 to-blue-600 text-white">
        <div className="px-4 mx-auto max-w-7xl lg:px-8 text-center">
          <h2 className="text-4xl font-bold mb-4">
            Ready to Start Your Journey?
          </h2>
          <p className="text-xl mb-8 text-purple-100 max-w-2xl mx-auto">
            Join thousands of successful students and take the first step towards your American citizenship.
          </p>
          <button 
            onClick={() => setCurrentView('courses')}
            className="px-8 py-4 text-lg font-semibold text-purple-600 bg-white rounded-xl hover:bg-gray-100 transition-all duration-300 transform hover:scale-105 shadow-lg hover:shadow-xl"
          >
            Get Started Today
          </button>
        </div>
      </section>
    </div>
  );

  const CoursesPage = () => (
    <div className="min-h-screen bg-gradient-to-b from-purple-100 to-blue-100 py-12">
      <div className="px-4 mx-auto max-w-7xl lg:px-8">
        <div className="text-center mb-12">
          <h1 className="text-4xl font-bold text-gray-900 mb-4">
            Citizenship Test Courses
          </h1>
          <p className="text-xl text-gray-700 max-w-3xl mx-auto">
            Choose from our expertly designed courses to prepare for your U.S. citizenship test
          </p>
        </div>
        
        <div className="grid gap-8 md:grid-cols-2 lg:grid-cols-3">
          {courses.map((course) => (
            <div key={course.id} className="bg-white/90 backdrop-blur-sm rounded-2xl shadow-lg hover:shadow-xl transition-all duration-300 overflow-hidden group border border-purple-100">
              <div className="relative">
                <img 
                  src={course.image}
                  alt={course.title}
                  className="w-full h-48 object-cover group-hover:scale-105 transition-transform duration-300"
                />
                <div className="absolute top-4 left-4 bg-white px-3 py-1 rounded-full text-sm font-medium text-purple-600">
                  {course.level}
                </div>
              </div>
              
              <div className="p-6">
                <h3 className="text-xl font-semibold text-gray-900 mb-2">{course.title}</h3>
                <p className="text-gray-600 mb-4 line-clamp-2">{course.description}</p>
                
                <div className="flex items-center gap-4 mb-4 text-sm text-gray-500">
                  <div className="flex items-center gap-1">
                    <Star className="w-4 h-4 text-yellow-400 fill-current" />
                    <span>{course.rating}</span>
                  </div>
                  <div className="flex items-center gap-1">
                    <Users className="w-4 h-4" />
                    <span>{course.students.toLocaleString()}</span>
                  </div>
                  <div className="flex items-center gap-1">
                    <Clock className="w-4 h-4" />
                    <span>{course.duration}</span>
                  </div>
                </div>
                
                <div className="mb-4">
                  <p className="text-sm text-gray-500 mb-2">Course modules:</p>
                  <div className="flex flex-wrap gap-2">
                    {course.modules.slice(0, 2).map((module, index) => (
                      <span key={index} className="px-2 py-1 bg-purple-100 text-purple-700 text-xs rounded-full">
                        {module}
                      </span>
                    ))}
                    {course.modules.length > 2 && (
                      <span className="px-2 py-1 bg-gray-100 text-gray-600 text-xs rounded-full">
                        +{course.modules.length - 2} more
                      </span>
                    )}
                  </div>
                </div>
                
                <div className="flex items-center justify-between">
                  <div>
                    <span className="text-2xl font-bold text-gray-900">${course.price}</span>
                  </div>
                  <button className="px-6 py-2 bg-gradient-to-r from-purple-600 to-blue-600 text-white rounded-lg hover:from-purple-700 hover:to-blue-700 transition-all duration-300 font-medium">
                    Enroll Now
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );

  const PricingPage = () => (
    <div className="min-h-screen bg-gradient-to-b from-purple-100 to-blue-100 py-12">
      <div className="px-4 mx-auto max-w-7xl lg:px-8">
        <div className="text-center mb-12">
          <h1 className="text-4xl font-bold text-gray-900 mb-4">
            Choose Your Learning Plan
          </h1>
          <p className="text-xl text-gray-700 max-w-3xl mx-auto">
            Select the perfect plan for your citizenship test preparation journey
          </p>
        </div>
        
        <div className="grid gap-8 md:grid-cols-3 max-w-5xl mx-auto">
          {pricingPlans.map((plan, index) => (
            <div key={index} className={`bg-white/90 backdrop-blur-sm rounded-2xl shadow-lg hover:shadow-2xl transition-all duration-300 overflow-hidden relative border ${plan.recommended ? 'ring-4 ring-purple-500 ring-opacity-50 scale-105 border-purple-200' : 'border-purple-100'}`}>
              {plan.recommended && (
                <div className="absolute top-0 left-0 right-0 bg-gradient-to-r from-purple-600 to-blue-600 text-white text-center py-2 text-sm font-semibold">
                  Most Popular
                </div>
              )}
              
              <div className={`p-8 ${plan.recommended ? 'pt-12' : ''}`}>
                <div className="text-center mb-8">
                  <div className="w-16 h-16 bg-gradient-to-r from-purple-600 to-blue-600 rounded-xl flex items-center justify-center mx-auto mb-4">
                    <plan.icon className="w-8 h-8 text-white" />
                  </div>
                  <h3 className="text-2xl font-bold text-gray-900 mb-2">{plan.name}</h3>
                  <div className="mb-4">
                    <span className="text-4xl font-bold text-gray-900">${plan.price}</span>
                    <span className="text-gray-600">/{plan.period}</span>
                  </div>
                </div>
                
                <ul className="space-y-4 mb-8">
                  {plan.features.map((feature, featureIndex) => (
                    <li key={featureIndex} className="flex items-center gap-3">
                      <CheckCircle className="w-5 h-5 text-green-500 flex-shrink-0" />
                      <span className="text-gray-700">{feature}</span>
                    </li>
                  ))}
                </ul>
                
                <button className={`w-full py-3 px-6 rounded-lg font-semibold transition-all duration-300 ${
                  plan.recommended 
                    ? 'bg-gradient-to-r from-purple-600 to-blue-600 text-white hover:from-purple-700 hover:to-blue-700 shadow-lg hover:shadow-xl' 
                    : 'bg-gray-100 text-gray-900 hover:bg-gray-200'
                }`}>
                  Get Started
                </button>
              </div>
            </div>
          ))}
        </div>
        
        <div className="text-center mt-12">
          <p className="text-gray-700 mb-4">
            All plans include a 30-day money-back guarantee
          </p>
          <div className="text-sm text-gray-600">
            * Premium plan includes exam fee refund if you don't pass within 6 months
          </div>
        </div>
      </div>
    </div>
  );

  const ContactPage = () => (
    <div className="min-h-screen bg-gradient-to-b from-purple-100 to-blue-100 py-12">
      <div className="px-4 mx-auto max-w-7xl lg:px-8">
        <div className="text-center mb-12">
          <h1 className="text-4xl font-bold text-gray-900 mb-4">
            Get in Touch
          </h1>
          <p className="text-xl text-gray-700 max-w-3xl mx-auto">
            Have questions about our courses? We're here to help you succeed in your citizenship journey.
          </p>
        </div>
        
        <div className="grid gap-12 lg:grid-cols-2 max-w-6xl mx-auto">
          {/* Contact Form */}
          <div className="bg-white/90 backdrop-blur-sm rounded-2xl shadow-lg p-8 border border-purple-100">
            <h2 className="text-2xl font-bold text-gray-900 mb-6">Send us a message</h2>
            <form className="space-y-6">
              <div className="grid gap-6 md:grid-cols-2">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    First Name
                  </label>
                  <input 
                    type="text" 
                    className="w-full px-4 py-3 border border-purple-200 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500 transition-colors bg-white/80"
                    placeholder="John"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Last Name
                  </label>
                  <input 
                    type="text" 
                    className="w-full px-4 py-3 border border-purple-200 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500 transition-colors bg-white/80"
                    placeholder="Doe"
                  />
                </div>
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Email
                </label>
                <input 
                  type="email" 
                  className="w-full px-4 py-3 border border-purple-200 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500 transition-colors bg-white/80"
                  placeholder="john@example.com"
                />
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Subject
                </label>
                <select className="w-full px-4 py-3 border border-purple-200 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500 transition-colors bg-white/80">
                  <option>General Inquiry</option>
                  <option>Course Questions</option>
                  <option>Technical Support</option>
                  <option>Billing</option>
                </select>
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Message
                </label>
                <textarea 
                  rows={6}
                  className="w-full px-4 py-3 border border-purple-200 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500 transition-colors bg-white/80"
                  placeholder="How can we help you?"
                ></textarea>
              </div>
              
              <button className="w-full bg-gradient-to-r from-purple-600 to-blue-600 text-white py-3 px-6 rounded-lg font-semibold hover:from-purple-700 hover:to-blue-700 transition-all duration-300 flex items-center justify-center gap-2">
                <Send className="w-5 h-5" />
                Send Message
              </button>
            </form>
          </div>
          
          {/* Contact Info */}
          <div className="space-y-8">
            <div className="bg-white/90 backdrop-blur-sm rounded-2xl shadow-lg p-8 border border-purple-100">
              <h2 className="text-2xl font-bold text-gray-900 mb-6">Contact Information</h2>
              <div className="space-y-6">
                <div className="flex items-start gap-4">
                  <div className="w-12 h-12 bg-purple-100 rounded-lg flex items-center justify-center flex-shrink-0">
                    <Mail className="w-6 h-6 text-purple-600" />
                  </div>
                  <div>
                    <h3 className="font-semibold text-gray-900">Email</h3>
                    <p className="text-gray-700">support@citizenshipprep.com</p>
                  </div>
                </div>
                
                <div className="flex items-start gap-4">
                  <div className="w-12 h-12 bg-purple-100 rounded-lg flex items-center justify-center flex-shrink-0">
                    <Phone className="w-6 h-6 text-purple-600" />
                  </div>
                  <div>
                    <h3 className="font-semibold text-gray-900">Phone</h3>
                    <p className="text-gray-700">1-800-CITIZEN</p>
                  </div>
                </div>
                
                <div className="flex items-start gap-4">
                  <div className="w-12 h-12 bg-purple-100 rounded-lg flex items-center justify-center flex-shrink-0">
                    <MapPin className="w-6 h-6 text-purple-600" />
                  </div>
                  <div>
                    <h3 className="font-semibold text-gray-900">Address</h3>
                    <p className="text-gray-700">
                      123 Learning Ave<br />
                      San Francisco, CA 94102
                    </p>
                  </div>
                </div>
              </div>
            </div>
            
            <div className="bg-gradient-to-r from-purple-600 to-blue-600 rounded-2xl p-8 text-white">
              <h3 className="text-xl font-bold mb-4">Need Immediate Help?</h3>
              <p className="mb-6 opacity-90">
                Join our live chat support for instant answers to your questions.
              </p>
              <button className="bg-white text-purple-600 px-6 py-3 rounded-lg font-semibold hover:bg-gray-100 transition-colors">
                Start Live Chat
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );

  const renderCurrentView = () => {
    switch (currentView) {
      case 'home':
        return <HomePage />;
      case 'courses':
        return <CoursesPage />;
      case 'pricing':
        return <PricingPage />;
      case 'contact':
        return <ContactPage />;
      default:
        return <HomePage />;
    }
  };

  return (
    <div className="min-h-screen bg-white">
      {/* Navigation */}
      <nav className="bg-white/90 backdrop-blur-md shadow-lg border-b border-purple-200 sticky top-0 z-50">
        <div className="px-4 mx-auto max-w-7xl lg:px-8">
          <div className="flex items-center justify-between h-16">
            {/* Logo */}
            <div className="flex items-center space-x-3">
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
                  onClick={() => setCurrentView(item.id)}
                  className={`flex items-center space-x-2 px-3 py-2 rounded-lg transition-colors ${
                    currentView === item.id 
                      ? 'text-purple-600 bg-purple-50' 
                      : 'text-gray-600 hover:text-gray-900 hover:bg-gray-50'
                  }`}
                >
                  <item.icon className="w-4 h-4" />
                  <span>{item.label}</span>
                </button>
              ))}
            </div>
            
            {/* Auth Buttons */}
            <div className="hidden md:flex items-center space-x-4">
              {!isAuthenticated ? (
                <>
                  <button className="text-gray-600 hover:text-gray-900 font-medium">
                    Log In
                  </button>
                  <button className="bg-gradient-to-r from-purple-600 to-blue-600 text-white px-4 py-2 rounded-lg font-medium hover:from-purple-700 hover:to-blue-700 transition-all duration-300">
                    Sign Up
                  </button>
                </>
              ) : (
                <div className="flex items-center space-x-3">
                  <div className="w-8 h-8 bg-gray-300 rounded-full flex items-center justify-center">
                    <User className="w-5 h-5 text-gray-600" />
                  </div>
                  <span className="text-gray-700 font-medium">
                    {userRole === 'student' ? 'Student' : 'Teacher'}
                  </span>
                </div>
              )}
            </div>
            
            {/* Mobile Menu Button */}
            <button
              className="md:hidden p-2"
              onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
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
            <div className="md:hidden py-4 border-t border-purple-200">
              <div className="flex flex-col space-y-2">
                {navItems.map((item) => (
                  <button
                    key={item.id}
                    onClick={() => {
                      setCurrentView(item.id);
                      setIsMobileMenuOpen(false);
                    }}
                    className={`flex items-center space-x-2 px-3 py-2 rounded-lg text-left ${
                      currentView === item.id 
                        ? 'text-purple-600 bg-purple-50' 
                        : 'text-gray-600'
                    }`}
                  >
                    <item.icon className="w-4 h-4" />
                    <span>{item.label}</span>
                  </button>
                ))}
                
                <div className="pt-4 border-t border-purple-200 space-y-2">
                  {!isAuthenticated ? (
                    <>
                      <button className="w-full text-left px-3 py-2 text-gray-600">
                        Log In
                      </button>
                      <button className="w-full text-left px-3 py-2 bg-gradient-to-r from-purple-600 to-blue-600 text-white rounded-lg">
                        Sign Up
                      </button>
                    </>
                  ) : (
                    <div className="px-3 py-2 flex items-center space-x-3">
                      <div className="w-8 h-8 bg-gray-300 rounded-full flex items-center justify-center">
                        <User className="w-5 h-5 text-gray-600" />
                      </div>
                      <span className="text-gray-700 font-medium">
                        {userRole === 'student' ? 'Student' : 'Teacher'}
                      </span>
                    </div>
                  )}
                </div>
              </div>
            </div>
          )}
        </div>
      </nav>

      {/* Main Content */}
      {renderCurrentView()}

      {/* Footer */}
      <footer className="bg-gray-900 text-white">
        <div className="px-4 py-16 mx-auto max-w-7xl lg:px-8">
          <div className="grid gap-8 md:grid-cols-2 lg:grid-cols-4">
            <div>
              <div className="flex items-center space-x-3 mb-6">
                <div className="w-10 h-10 bg-gradient-to-r from-purple-600 to-blue-600 rounded-xl flex items-center justify-center">
                  <BookOpen className="w-6 h-6 text-white" />
                </div>
                <span className="text-xl font-bold">CitizenshipPrep</span>
              </div>
              <p className="text-gray-400 leading-relaxed">
                Helping aspiring citizens achieve their American dream through comprehensive test preparation and expert guidance.
              </p>
            </div>
            
            <div>
              <h3 className="text-lg font-semibold mb-4">Courses</h3>
              <ul className="space-y-2 text-gray-400">
                <li><a href="#" className="hover:text-white transition-colors">Citizenship Test Prep</a></li>
                <li><a href="#" className="hover:text-white transition-colors">English for Citizenship</a></li>
                <li><a href="#" className="hover:text-white transition-colors">Mock Interviews</a></li>
                <li><a href="#" className="hover:text-white transition-colors">Practice Tests</a></li>
              </ul>
            </div>
            
            <div>
              <h3 className="text-lg font-semibold mb-4">Support</h3>
              <ul className="space-y-2 text-gray-400">
                <li><a href="#" className="hover:text-white transition-colors">Help Center</a></li>
                <li><a href="#" className="hover:text-white transition-colors">Live Chat</a></li>
                <li><a href="#" className="hover:text-white transition-colors">Contact Us</a></li>
                <li><a href="#" className="hover:text-white transition-colors">FAQ</a></li>
              </ul>
            </div>
            
            <div>
              <h3 className="text-lg font-semibold mb-4">Company</h3>
              <ul className="space-y-2 text-gray-400">
                <li><a href="#" className="hover:text-white transition-colors">About Us</a></li>
                <li><a href="#" className="hover:text-white transition-colors">Privacy Policy</a></li>
                <li><a href="#" className="hover:text-white transition-colors">Terms of Service</a></li>
                <li><a href="#" className="hover:text-white transition-colors">Blog</a></li>
              </ul>
            </div>
          </div>
          
          <div className="pt-8 mt-8 border-t border-gray-800 text-center text-gray-400">
            <p>&copy; 2025 CitizenshipPrep. All rights reserved.</p>
          </div>
        </div>
      </footer>
    </div>
  );
};

export default LearningManagementSystem;
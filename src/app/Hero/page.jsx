"use client"
import React from 'react';
import { 
  BookOpen, 
  Users, 
  Award, 
  Star, 
  CheckCircle
} from 'lucide-react';

const Heropage = ({ setCurrentView }) => {
  return (
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
};

export default Heropage;
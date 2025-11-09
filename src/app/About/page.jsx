import React from 'react';
import { ArrowRight, BookOpen, Globe, Users, Award, Target, Heart, CheckCircle } from 'lucide-react';

const AboutPage = () => {
  return (
    <div className="min-h-screen bg-gradient-to-b from-purple-100 to-blue-100">
      <div className="px-4 py-12 mx-auto max-w-7xl lg:px-8">
        {/* Header */}
        <header className="mb-16 text-center">
          <h1 className="text-5xl font-bold text-gray-900 mb-6">About CitizenshipPrep</h1>
          <p className="text-xl text-gray-700 max-w-3xl mx-auto leading-relaxed">
            Empowering your journey to U.S. citizenship through expert instruction, 
            innovative technology, and personalized support.
          </p>
        </header>

        {/* Mission Section */}
        <section className="bg-white/90 backdrop-blur-sm rounded-2xl shadow-lg p-8 mb-12 border border-purple-100">
          <div className="grid gap-8 lg:grid-cols-2 items-center">
            <div>
              <div className="flex items-center gap-3 mb-6">
                <div className="w-12 h-12 bg-gradient-to-r from-purple-600 to-blue-600 rounded-xl flex items-center justify-center">
                  <Target className="w-6 h-6 text-white" />
                </div>
                <h2 className="text-3xl font-bold text-gray-900">Our Mission</h2>
              </div>
              <p className="text-gray-700 mb-6 leading-relaxed">
                At CitizenshipPrep, our goal is to help people learn English effectively, 
                focusing on the skills needed to pass the citizenship test. We specialize in 
                assisting Spanish and Chinese speakers in their journey to becoming 
                confident English communicators and successful U.S. citizens.
              </p>
              <p className="text-gray-700 leading-relaxed">
                We believe that language should never be a barrier to achieving your dreams. 
                Our dedicated team works tirelessly to create engaging, accessible, and 
                effective learning materials tailored to your unique needs and background.
              </p>
            </div>
            <div className="relative">
              <div className="absolute inset-0 bg-gradient-to-r from-purple-400 to-blue-500 rounded-2xl rotate-3 opacity-20"></div>
              <img 
                src="https://images.unsplash.com/photo-1522202176988-66273c2fd55f?w=500&h=400&fit=crop"
                alt="Students learning together"
                className="relative z-10 w-full h-80 object-cover rounded-2xl shadow-lg"
              />
            </div>
          </div>
        </section>

        {/* What We Offer */}
        <section className="mb-16">
          <div className="text-center mb-12">
            <h2 className="text-4xl font-bold text-gray-900 mb-4">What We Offer</h2>
            <p className="text-xl text-gray-700 max-w-3xl mx-auto">
              Comprehensive solutions designed to support your citizenship journey every step of the way
            </p>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            {[
              { 
                title: 'Tailored Curriculum', 
                description: 'Courses designed specifically for the citizenship test with culturally relevant examples',
                icon: <BookOpen className="w-6 h-6 text-purple-600" />
              },
              { 
                title: 'Language Support', 
                description: 'Specialized help for Spanish and Chinese speakers with native language explanations',
                icon: <Globe className="w-6 h-6 text-blue-600" />
              },
              { 
                title: 'Interactive Learning', 
                description: 'Engaging exercises, simulations, and practice tests to build your confidence',
                icon: <Users className="w-6 h-6 text-purple-600" />
              },
              { 
                title: 'Progress Tracking', 
                description: 'Advanced analytics to monitor your improvement and identify areas for focus',
                icon: <ArrowRight className="w-6 h-6 text-blue-600" />
              }
            ].map((feature, index) => (
              <div key={index} className="bg-white/90 backdrop-blur-sm rounded-2xl shadow-lg p-6 border border-purple-100 hover:shadow-xl transition-all duration-300 group">
                <div className="w-14 h-14 bg-gradient-to-r from-purple-100 to-blue-100 rounded-xl flex items-center justify-center mb-4 group-hover:scale-110 transition-transform duration-300">
                  {feature.icon}
                </div>
                <h3 className="text-lg font-semibold text-gray-900 mb-3">{feature.title}</h3>
                <p className="text-gray-700 leading-relaxed">{feature.description}</p>
              </div>
            ))}
          </div>
        </section>

        {/* Our Values */}
        <section className="mb-16">
          <div className="bg-white/90 backdrop-blur-sm rounded-2xl shadow-lg p-8 border border-purple-100">
            <div className="text-center mb-12">
              <div className="flex items-center justify-center gap-3 mb-6">
                <div className="w-12 h-12 bg-gradient-to-r from-purple-600 to-blue-600 rounded-xl flex items-center justify-center">
                  <Heart className="w-6 h-6 text-white" />
                </div>
                <h2 className="text-3xl font-bold text-gray-900">Our Values</h2>
              </div>
              <p className="text-xl text-gray-700 max-w-2xl mx-auto">
                The principles that guide everything we do
              </p>
            </div>
            
            <div className="grid gap-8 md:grid-cols-3">
              {[
                {
                  title: 'Accessibility',
                  description: 'Education should be available to everyone, regardless of their background or starting point.'
                },
                {
                  title: 'Excellence',
                  description: 'We maintain the highest standards in our content, technology, and student support.'
                },
                {
                  title: 'Empowerment',
                  description: 'We believe in empowering students with the tools and confidence to achieve their dreams.'
                }
              ].map((value, index) => (
                <div key={index} className="text-center">
                  <div className="w-16 h-16 bg-gradient-to-r from-purple-600 to-blue-600 rounded-full flex items-center justify-center mx-auto mb-4">
                    <CheckCircle className="w-8 h-8 text-white" />
                  </div>
                  <h3 className="text-xl font-semibold text-gray-900 mb-3">{value.title}</h3>
                  <p className="text-gray-700 leading-relaxed">{value.description}</p>
                </div>
              ))}
            </div>
          </div>
        </section>

        {/* Team Section */}
        <section className="mb-16">
          <div className="text-center mb-12">
            <h2 className="text-4xl font-bold text-gray-900 mb-4">Meet Our Team</h2>
            <p className="text-xl text-gray-700 max-w-3xl mx-auto">
              Experienced educators and immigration experts dedicated to your success
            </p>
          </div>
          
          <div className="grid gap-8 md:grid-cols-3">
            {[
              {
                name: 'Sarah Johnson',
                role: 'Lead Citizenship Instructor',
                image: 'https://images.unsplash.com/photo-1494790108755-2616b612b1e7?w=300&h=300&fit=crop&crop=face',
                description: '10+ years helping students pass citizenship tests'
              },
              {
                name: 'Michael Chen',
                role: 'ESL Specialist',
                image: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=300&h=300&fit=crop&crop=face',
                description: 'Expert in English learning for Chinese speakers'
              },
              {
                name: 'Emily Rodriguez',
                role: 'Immigration Counselor',
                image: 'https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=300&h=300&fit=crop&crop=face',
                description: 'Certified immigration consultant and interview coach'
              }
            ].map((member, index) => (
              <div key={index} className="bg-white/90 backdrop-blur-sm rounded-2xl shadow-lg p-6 border border-purple-100 text-center hover:shadow-xl transition-all duration-300">
                <img 
                  src={member.image}
                  alt={member.name}
                  className="w-24 h-24 rounded-full mx-auto mb-4 object-cover border-4 border-purple-100"
                />
                <h3 className="text-xl font-semibold text-gray-900 mb-2">{member.name}</h3>
                <p className="text-purple-600 font-medium mb-3">{member.role}</p>
                <p className="text-gray-700 text-sm">{member.description}</p>
              </div>
            ))}
          </div>
        </section>

        {/* CTA Section */}
        <section className="bg-gradient-to-r from-purple-600 to-blue-600 rounded-2xl p-12 text-white text-center">
          <div className="max-w-3xl mx-auto">
            <div className="w-16 h-16 bg-white/20 rounded-full flex items-center justify-center mx-auto mb-6">
              <Award className="w-8 h-8 text-white" />
            </div>
            <h2 className="text-4xl font-bold mb-6">Start Your Journey Today</h2>
            <p className="text-xl mb-8 text-purple-100 leading-relaxed">
              Join thousands of successful learners who have achieved their dreams of U.S. citizenship. 
              Take the first step towards your new life as an American citizen.
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <button className="bg-white text-purple-600 font-bold py-3 px-8 rounded-xl transition-all duration-300 hover:bg-gray-100 transform hover:scale-105 shadow-lg hover:shadow-xl">
                Begin Learning
              </button>
              <button className="border-2 border-white text-white font-bold py-3 px-8 rounded-xl transition-all duration-300 hover:bg-white hover:text-purple-600 transform hover:scale-105">
                Schedule Consultation
              </button>
            </div>
          </div>
        </section>
      </div>
    </div>
  );
};

export default AboutPage;
import React from 'react';
import { ArrowRight, BookOpen, Globe, Users, Award, Target, Heart, CheckCircle } from 'lucide-react';

const AboutPage = () => {

  const features = [
    { 
      title: 'Citizenship Civics Learning', 
      description: 'Clear, simple lessons covering all 100 official citizenship questions, organized by topic for easy study.',
      icon: <BookOpen className="w-6 h-6 text-red-600" />
    },
    { 
      title: 'Citizenship English', 
      description: 'English practice focused on the words, phrases, and questions used in real citizenship interviews.',
      icon: <Globe className="w-6 h-6 text-blue-600" />
    },
    { 
      title: 'Practice Tests & Quizzes', 
      description: 'Interactive quizzes and full practice exams that help you test your knowledge and build confidence.',
      icon: <Users className="w-6 h-6 text-red-600" />
    },
    { 
      title: 'Readiness Tracking', 
      description: 'A personalized dashboard that shows your progress, strengths, and how close you are to being ready.',
      icon: <ArrowRight className="w-6 h-6 text-blue-600" />
    }
  ]

  const values = [
    {
      title: 'Accessibility',
      description: 'We believe everyone deserves access to high-quality learning, regardless of language, education level, or background.'
    },
    {
      title: 'Clarity',
      description: 'We focus on simple explanations, clear steps, and practical guidance so learners always know what to do next.'
    },
    {
      title: 'Empowerment',
      description: 'Our goal is not just to teach facts, but to help people feel confident, prepared, and in control of their future.'
    }
  ]

  const teamMembers = [
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
  ]


  return (
    <div className="min-h-screen ">
      <div className="px-4 py-12 mx-auto max-w-7xl lg:px-8">
        {/* Header */}
        <header className="mb-16 text-center">
          <h1 className="text-5xl font-bold text-gray-900 mb-6">About CitizenshipPrep</h1>
          <p className="text-xl text-gray-700 max-w-3xl mx-auto leading-relaxed">
            A digital companion built to help immigrants prepare for U.S. citizenship, build confidence in English, and navigate life in the United States.
          </p>
        </header>

        {/* Mission Section */}
        <section className="bg-white/90 backdrop-blur-sm rounded-2xl shadow-lg p-8 mb-12 border border-red-100">
          <div className="grid gap-8 lg:grid-cols-2 items-center">
            <div>
              <div className="flex items-center gap-3 mb-6">
                <div className="w-12 h-12 bg-gradient-to-r from-red-600 to-blue-600 rounded-xl flex items-center justify-center">
                  <Target className="w-6 h-6 text-white" />
                </div>
                <h2 className="text-3xl font-bold text-gray-900">Our Mission</h2>
              </div>
              <p className="text-gray-700 mb-6 leading-relaxed">
                CitizenshipPrep exists to give immigrants a clear, supportive path toward U.S. citizenship. 
                Instead of studying randomly, learners can follow structured lessons, take realistic practice 
                tests, and see exactly how ready they are for their citizenship interview.
              </p>
              <p className="text-gray-700 leading-relaxed">
                We combine simple English, civics education, and smart progress tracking so that every learner 
                can understand what to study, what to practice next, and when they are ready to move forward with confidence.
              </p>
            </div>
            <div className="relative">
              <div className="absolute inset-0 bg-gradient-to-r from-red-400 to-blue-500 rounded-2xl rotate-3 opacity-20"></div>
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
            {features.map((feature, index) => (
              <div key={index} className="bg-white/90 backdrop-blur-sm rounded-2xl shadow-lg p-6 border border-red-100 hover:shadow-xl transition-all duration-300 group">
                <div className="w-14 h-14 bg-gradient-to-r from-red-100 to-blue-100 rounded-xl flex items-center justify-center mb-4 group-hover:scale-110 transition-transform duration-300">
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
          <div className="bg-white/90 backdrop-blur-sm rounded-2xl shadow-lg p-8 border border-red-100">
            <div className="text-center mb-12">
              <div className="flex items-center justify-center gap-3 mb-6">
                <div className="w-12 h-12 bg-gradient-to-r from-red-600 to-blue-600 rounded-xl flex items-center justify-center">
                  <Heart className="w-6 h-6 text-white" />
                </div>
                <h2 className="text-3xl font-bold text-gray-900">Our Values</h2>
              </div>
              <p className="text-xl text-gray-700 max-w-2xl mx-auto">
                The principles that guide everything we do
              </p>
            </div>
            
            <div className="grid gap-8 md:grid-cols-3">
              {values.map((value, index) => (
                <div key={index} className="text-center">
                  <div className="w-16 h-16 bg-gradient-to-r from-red-600 to-blue-600 rounded-full flex items-center justify-center mx-auto mb-4">
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
              A growing team of technologists, educators, and immigrant advocates building tools to make the citizenship journey easier.
            </p>
          </div>
          
          <div className="grid gap-8 md:grid-cols-3">
            {teamMembers.map((member, index) => (
              <div key={index} className="bg-white/90 backdrop-blur-sm rounded-2xl shadow-lg p-6 border border-red-100 text-center hover:shadow-xl transition-all duration-300">
                <img 
                  src={member.image}
                  alt={member.name}
                  className="w-24 h-24 rounded-full mx-auto mb-4 object-cover border-4 border-red-100"
                />
                <h3 className="text-xl font-semibold text-gray-900 mb-2">{member.name}</h3>
                <p className="text-red-600 font-medium mb-3">{member.role}</p>
                <p className="text-gray-700 text-sm">{member.description}</p>
              </div>
            ))}
          </div>
        </section>

        
      </div>
    </div>
  );
};

export default AboutPage;
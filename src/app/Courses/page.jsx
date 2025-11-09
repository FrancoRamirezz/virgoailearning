import React from 'react';
import { Star, Users, Clock } from 'lucide-react';

const CoursesPage = () => {
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

  return (
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
};

export default CoursesPage;
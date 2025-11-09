"use client"
import React, { useState } from 'react';
import { 
  Play, 
  Pause, 
  SkipForward, 
  SkipBack, 
  FileText, 
  Download, 
  CheckCircle, 
  Clock,
  Users,
  Star,
  BookOpen,
  MessageSquare,
  Trophy,
  Target,
  Volume2
} from 'lucide-react';

const IndividualCourse = ({ courseId }) => {
  const [activeLesson, setActiveLesson] = useState(0);
  const [isPlaying, setIsPlaying] = useState(false);
  const [progress, setProgress] = useState(45); // Video progress percentage
  const [completedLessons, setCompletedLessons] = useState([0, 1, 2]);

  // Sample course data
  const course = {
    id: courseId,
    title: 'Citizenship Test Preparation',
    instructor: 'Sarah Johnson',
    description: 'Complete preparation for the U.S. citizenship test with interactive lessons and practice tests.',
    rating: 4.8,
    totalStudents: 2456,
    totalDuration: '12 hours',
    level: 'Beginner to Advanced',
    completion: 65,
    modules: [
      {
        id: 1,
        title: 'American History Fundamentals',
        lessons: [
          { id: 1, title: 'Colonial Period Overview', duration: '15:30', type: 'video', completed: true },
          { id: 2, title: 'Revolutionary War Key Events', duration: '18:45', type: 'video', completed: true },
          { id: 3, title: 'Constitution Formation', duration: '22:15', type: 'video', completed: true },
          { id: 4, title: 'History Quiz 1', duration: '10:00', type: 'quiz', completed: false }
        ]
      },
      {
        id: 2,
        title: 'Government Structure',
        lessons: [
          { id: 5, title: 'Three Branches of Government', duration: '20:30', type: 'video', completed: true },
          { id: 6, title: 'Federal vs State Powers', duration: '16:20', type: 'video', completed: false },
          { id: 7, title: 'Electoral Process', duration: '19:45', type: 'video', completed: false },
          { id: 8, title: 'Government Quiz', duration: '15:00', type: 'quiz', completed: false }
        ]
      },
      {
        id: 3,
        title: 'Civics and Rights',
        lessons: [
          { id: 9, title: 'Bill of Rights Deep Dive', duration: '25:15', type: 'video', completed: false },
          { id: 10, title: 'Citizenship Rights & Responsibilities', duration: '17:30', type: 'video', completed: false },
          { id: 11, title: 'Legal System Overview', duration: '21:00', type: 'video', completed: false },
          { id: 12, title: 'Final Practice Test', duration: '30:00', type: 'test', completed: false }
        ]
      }
    ]
  };

  const currentLesson = course.modules
    .flatMap(module => module.lessons)
    .find((_, index) => index === activeLesson) || course.modules[0].lessons[0];

  const totalLessons = course.modules.reduce((total, module) => total + module.lessons.length, 0);
  const completedCount = course.modules.reduce((total, module) => 
    total + module.lessons.filter(lesson => lesson.completed).length, 0);

  const handleLessonComplete = () => {
    setCompletedLessons([...completedLessons, activeLesson]);
    // Move to next lesson automatically
    if (activeLesson < totalLessons - 1) {
      setActiveLesson(activeLesson + 1);
    }
  };

  const getLessonIcon = (type) => {
    switch (type) {
      case 'video': return Play;
      case 'quiz': return FileText;
      case 'test': return Target;
      default: return BookOpen;
    }
  };

  const getLessonColor = (type) => {
    switch (type) {
      case 'video': return 'text-blue-600 bg-blue-100';
      case 'quiz': return 'text-purple-600 bg-purple-100';
      case 'test': return 'text-red-600 bg-red-100';
      default: return 'text-gray-600 bg-gray-100';
    }
  };

  return (
    <div className="grid grid-cols-1 lg:grid-cols-4 gap-6 h-screen overflow-hidden">
      {/* Sidebar - Course Content */}
      <div className="lg:col-span-1 bg-white/90 backdrop-blur-sm rounded-2xl shadow-lg border border-purple-100 p-6 overflow-y-auto">
        <div className="mb-6">
          <h3 className="text-lg font-bold text-gray-900 mb-2">{course.title}</h3>
          <p className="text-sm text-gray-600 mb-4">by {course.instructor}</p>
          
          {/* Progress */}
          <div className="mb-4">
            <div className="flex justify-between text-sm text-gray-600 mb-2">
              <span>Progress</span>
              <span>{completedCount}/{totalLessons} lessons</span>
            </div>
            <div className="w-full bg-gray-200 rounded-full h-2">
              <div
                className="h-2 bg-gradient-to-r from-purple-600 to-blue-600 rounded-full transition-all duration-500"
                style={{ width: `${(completedCount / totalLessons) * 100}%` }}
              ></div>
            </div>
          </div>
        </div>

        {/* Course Modules */}
        <div className="space-y-4">
          {course.modules.map((module) => (
            <div key={module.id} className="border border-gray-200 rounded-xl overflow-hidden">
              <div className="p-4 bg-gradient-to-r from-purple-50 to-blue-50 border-b border-gray-200">
                <h4 className="font-semibold text-gray-900">{module.title}</h4>
              </div>
              
              <div className="divide-y divide-gray-200">
                {module.lessons.map((lesson, lessonIndex) => {
                  const globalIndex = course.modules
                    .slice(0, module.id - 1)
                    .reduce((total, mod) => total + mod.lessons.length, 0) + lessonIndex;
                  
                  const LessonIcon = getLessonIcon(lesson.type);
                  const isActive = globalIndex === activeLesson;
                  
                  return (
                    <button
                      key={lesson.id}
                      onClick={() => setActiveLesson(globalIndex)}
                      className={`w-full p-3 text-left hover:bg-gray-50 transition-colors duration-200 ${
                        isActive ? 'bg-purple-50 border-r-4 border-purple-600' : ''
                      }`}
                    >
                      <div className="flex items-center space-x-3">
                        <div className={`w-8 h-8 rounded-lg ${getLessonColor(lesson.type)} flex items-center justify-center`}>
                          {lesson.completed ? (
                            <CheckCircle className="w-4 h-4 text-green-600" />
                          ) : (
                            <LessonIcon className="w-4 h-4" />
                          )}
                        </div>
                        
                        <div className="flex-1 min-w-0">
                          <p className={`text-sm font-medium truncate ${
                            isActive ? 'text-purple-700' : 'text-gray-900'
                          }`}>
                            {lesson.title}
                          </p>
                          <div className="flex items-center space-x-2 text-xs text-gray-500">
                            <Clock className="w-3 h-3" />
                            <span>{lesson.duration}</span>
                            <span className="capitalize">• {lesson.type}</span>
                          </div>
                        </div>
                      </div>
                    </button>
                  );
                })}
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Main Content Area */}
      <div className="lg:col-span-3 space-y-6">
        {/* Video Player */}
        <div className="bg-white/90 backdrop-blur-sm rounded-2xl shadow-lg border border-purple-100 overflow-hidden">
          <div className="relative bg-gray-900 aspect-video">
            {/* Video placeholder */}
            <div className="absolute inset-0 flex items-center justify-center">
              <div className="text-center text-white">
                <div className="w-20 h-20 bg-white/20 rounded-full flex items-center justify-center mx-auto mb-4">
                  {isPlaying ? (
                    <Pause className="w-8 h-8" />
                  ) : (
                    <Play className="w-8 h-8 ml-1" />
                  )}
                </div>
                <h3 className="text-xl font-semibold mb-2">{currentLesson.title}</h3>
                <p className="text-gray-300">Duration: {currentLesson.duration}</p>
              </div>
            </div>
            
            {/* Video Controls Overlay */}
            <div className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/50 to-transparent p-4">
              <div className="flex items-center space-x-4">
                <button
                  onClick={() => setIsPlaying(!isPlaying)}
                  className="w-10 h-10 bg-white/20 backdrop-blur-sm rounded-full flex items-center justify-center hover:bg-white/30 transition-colors duration-200"
                >
                  {isPlaying ? (
                    <Pause className="w-5 h-5 text-white" />
                  ) : (
                    <Play className="w-5 h-5 text-white ml-0.5" />
                  )}
                </button>
                
                <button className="w-8 h-8 bg-white/20 backdrop-blur-sm rounded-full flex items-center justify-center hover:bg-white/30 transition-colors duration-200">
                  <SkipBack className="w-4 h-4 text-white" />
                </button>
                
                <button className="w-8 h-8 bg-white/20 backdrop-blur-sm rounded-full flex items-center justify-center hover:bg-white/30 transition-colors duration-200">
                  <SkipForward className="w-4 h-4 text-white" />
                </button>
                
                <div className="flex-1 mx-4">
                  <div className="w-full bg-white/20 rounded-full h-1">
                    <div
                      className="h-1 bg-purple-500 rounded-full transition-all duration-300"
                      style={{ width: `${progress}%` }}
                    ></div>
                  </div>
                </div>
                
                <button className="w-8 h-8 bg-white/20 backdrop-blur-sm rounded-full flex items-center justify-center hover:bg-white/30 transition-colors duration-200">
                  <Volume2 className="w-4 h-4 text-white" />
                </button>
              </div>
            </div>
          </div>
          
          {/* Lesson Info */}
          <div className="p-6">
            <div className="flex items-center justify-between mb-4">
              <div>
                <h2 className="text-2xl font-bold text-gray-900">{currentLesson.title}</h2>
                <p className="text-gray-600">Module: American History Fundamentals</p>
              </div>
              
              <div className="flex items-center space-x-3">
                <button className="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors duration-200 flex items-center gap-2">
                  <Download className="w-4 h-4" />
                  Download
                </button>
                
                <button
                  onClick={handleLessonComplete}
                  className="px-6 py-2 bg-gradient-to-r from-purple-600 to-blue-600 text-white rounded-lg hover:from-purple-700 hover:to-blue-700 transition-all duration-200 flex items-center gap-2"
                >
                  <CheckCircle className="w-4 h-4" />
                  Mark Complete
                </button>
              </div>
            </div>
            
            {/* Lesson Description */}
            <div className="bg-gray-50 rounded-xl p-4">
              <p className="text-gray-700">
                In this lesson, we'll explore the colonial period of American history, covering key events, 
                major figures, and the foundations that led to the formation of the United States. 
                Understanding this period is crucial for the citizenship test.
              </p>
            </div>
          </div>
        </div>

        {/* Additional Resources & Discussion */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* Resources */}
          <div className="bg-white/90 backdrop-blur-sm rounded-2xl shadow-lg border border-purple-100 p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center gap-2">
              <FileText className="w-5 h-5 text-purple-600" />
              Additional Resources
            </h3>
            
            <div className="space-y-3">
              {[
                { title: 'Colonial Period Timeline PDF', type: 'PDF', size: '2.3 MB' },
                { title: 'Key Figures Reference Sheet', type: 'PDF', size: '1.8 MB' },
                { title: 'Practice Questions Set 1', type: 'Quiz', questions: '15 questions' }
              ].map((resource, index) => (
                <div key={index} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors duration-200 cursor-pointer">
                  <div className="flex items-center space-x-3">
                    <FileText className="w-4 h-4 text-gray-600" />
                    <div>
                      <p className="font-medium text-gray-900">{resource.title}</p>
                      <p className="text-xs text-gray-500">
                        {resource.size || resource.questions}
                      </p>
                    </div>
                  </div>
                  <Download className="w-4 h-4 text-gray-400" />
                </div>
              ))}
            </div>
          </div>

          {/* Discussion */}
          <div className="bg-white/90 backdrop-blur-sm rounded-2xl shadow-lg border border-purple-100 p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center gap-2">
              <MessageSquare className="w-5 h-5 text-blue-600" />
              Discussion
            </h3>
            
            <div className="space-y-3">
              {[
                { user: 'Maria S.', comment: 'Great explanation of the colonial period! The timeline really helped.', time: '2 hours ago' },
                { user: 'James K.', comment: 'Can someone clarify the difference between the French and Indian War impacts?', time: '5 hours ago' },
                { user: 'Sarah (Instructor)', comment: 'Remember to focus on the key events that directly influenced independence.', time: '1 day ago' }
              ].map((discussion, index) => (
                <div key={index} className="p-3 bg-gray-50 rounded-lg">
                  <div className="flex items-center justify-between mb-2">
                    <span className="font-medium text-gray-900">{discussion.user}</span>
                    <span className="text-xs text-gray-500">{discussion.time}</span>
                  </div>
                  <p className="text-sm text-gray-700">{discussion.comment}</p>
                </div>
              ))}
            </div>
            
            <div className="mt-4 pt-4 border-t border-gray-200">
              <div className="flex space-x-3">
                <input
                  type="text"
                  placeholder="Add to discussion..."
                  className="flex-1 px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500 text-sm"
                />
                <button className="px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors duration-200 text-sm">
                  Post
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default IndividualCourse;
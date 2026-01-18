"use client"
import React, { useState } from 'react';
import { ChevronLeft, ChevronRight, Plus, Clock, MapPin, Users, BookOpen, FileText } from 'lucide-react';

const Calendar = () => {
  const [currentDate, setCurrentDate] = useState(new Date());
  const [selectedDate, setSelectedDate] = useState(new Date());
  const [view, setView] = useState('month'); // 'month', 'week', 'day'
  const [showEventModal, setShowEventModal] = useState(false);

  // Sample events data
  const events = [
    {
      id: 1,
      title: 'Citizenship Test Practice',
      date: new Date(2025, 0, 15),
      time: '10:00 AM',
      duration: '2 hours',
      type: 'assignment',
      course: 'Citizenship Prep',
      location: 'Online',
      color: 'bg-purple-100 text-purple-800 border-purple-300'
    },
    {
      id: 2,
      title: 'Live Tutoring Session',
      date: new Date(2025, 0, 18),
      time: '2:00 PM',
      duration: '1 hour',
      type: 'meeting',
      course: 'English for Citizenship',
      location: 'Zoom',
      color: 'bg-blue-100 text-blue-800 border-blue-300'
    },
    {
      id: 3,
      title: 'Mock Interview',
      date: new Date(2025, 0, 20),
      time: '11:00 AM',
      duration: '30 minutes',
      type: 'exam',
      course: 'Interview Prep',
      location: 'Online',
      color: 'bg-red-100 text-red-800 border-red-300'
    },
    {
      id: 4,
      title: 'Assignment Due: History Quiz',
      date: new Date(2025, 0, 22),
      time: '11:59 PM',
      duration: 'Due',
      type: 'deadline',
      course: 'American History',
      location: 'Online Submission',
      color: 'bg-orange-100 text-orange-800 border-orange-300'
    }
  ];

  const months = [
    'January', 'February', 'March', 'April', 'May', 'June',
    'July', 'August', 'September', 'October', 'November', 'December'
  ];

  const daysOfWeek = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];

  const getDaysInMonth = (date) => {
    const year = date.getFullYear();
    const month = date.getMonth();
    const firstDay = new Date(year, month, 1);
    const lastDay = new Date(year, month + 1, 0);
    const daysInMonth = lastDay.getDate();
    const startingDayOfWeek = firstDay.getDay();

    const days = [];

    // Add empty cells for days before the first day of the month
    for (let i = 0; i < startingDayOfWeek; i++) {
      days.push(null);
    }

    // Add days of the month
    for (let day = 1; day <= daysInMonth; day++) {
      days.push(new Date(year, month, day));
    }

    return days;
  };

  const getEventsForDate = (date) => {
    if (!date) return [];
    return events.filter(event => 
      event.date.toDateString() === date.toDateString()
    );
  };

  const navigateMonth = (direction) => {
    const newDate = new Date(currentDate);
    newDate.setMonth(currentDate.getMonth() + direction);
    setCurrentDate(newDate);
  };

  const isToday = (date) => {
    if (!date) return false;
    const today = new Date();
    return date.toDateString() === today.toDateString();
  };

  const isSelected = (date) => {
    if (!date) return false;
    return date.toDateString() === selectedDate.toDateString();
  };

  const EventTypeIcon = ({ type }) => {
    switch (type) {
      case 'assignment':
        return <BookOpen className="w-4 h-4" />;
      case 'meeting':
        return <Users className="w-4 h-4" />;
      case 'exam':
        return <FileText className="w-4 h-4" />;
      case 'deadline':
        return <Clock className="w-4 h-4" />;
      default:
        return <Clock className="w-4 h-4" />;
    }
  };

  return (
    <div className="bg-white/90 backdrop-blur-sm rounded-2xl shadow-lg border border-purple-100 overflow-hidden">
      {/* Calendar Header */}
      <div className="p-6 border-b border-purple-100 bg-gradient-to-r from-purple-50 to-blue-50">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-2xl font-bold text-gray-900">Calendar</h2>
          <button 
            onClick={() => setShowEventModal(true)}
            className="px-4 py-2 bg-gradient-to-r from-purple-600 to-blue-600 text-white rounded-lg hover:from-purple-700 hover:to-blue-700 transition-all duration-200 flex items-center gap-2 font-medium"
          >
            <Plus className="w-4 h-4" />
            Add Event
          </button>
        </div>
        
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-4">
            <button
              onClick={() => navigateMonth(-1)}
              className="p-2 rounded-lg hover:bg-white/50 transition-colors duration-200"
              aria-label="Previous month"
            >
              <ChevronLeft className="w-5 h-5 text-gray-600" />
            </button>
            
            <h3 className="text-xl font-semibold text-gray-900">
              {months[currentDate.getMonth()]} {currentDate.getFullYear()}
            </h3>
            
            <button
              onClick={() => navigateMonth(1)}
              className="p-2 rounded-lg hover:bg-white/50 transition-colors duration-200"
              aria-label="Next month"
            >
              <ChevronRight className="w-5 h-5 text-gray-600" />
            </button>
          </div>
          
          <div className="flex items-center space-x-2">
            {['month', 'week', 'day'].map((viewType) => (
              <button
                key={viewType}
                onClick={() => setView(viewType)}
                className={`px-3 py-1 rounded-lg capitalize transition-colors duration-200 ${
                  view === viewType
                    ? 'bg-white text-purple-600 shadow-sm'
                    : 'text-gray-600 hover:bg-white/50'
                }`}
              >
                {viewType}
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* Calendar Grid */}
      <div className="p-6">
        {/* Days of week header */}
        <div className="grid grid-cols-7 gap-1 mb-4">
          {daysOfWeek.map((day) => (
            <div key={day} className="p-2 text-center text-sm font-medium text-gray-500">
              {day}
            </div>
          ))}
        </div>

        {/* Calendar days */}
        <div className="grid grid-cols-7 gap-1">
          {getDaysInMonth(currentDate).map((date, index) => {
            const dayEvents = getEventsForDate(date);
            const hasEvents = dayEvents.length > 0;

            return (
              <div
                key={index}
                className={`min-h-24 p-2 border rounded-lg cursor-pointer transition-all duration-200 ${
                  !date 
                    ? 'bg-gray-50' 
                    : isSelected(date)
                      ? 'bg-purple-100 border-purple-300 shadow-sm'
                      : isToday(date)
                        ? 'bg-blue-50 border-blue-300'
                        : hasEvents
                          ? 'bg-gray-50 border-gray-200 hover:bg-gray-100'
                          : 'bg-white border-gray-200 hover:bg-gray-50'
                }`}
                onClick={() => date && setSelectedDate(date)}
              >
                {date && (
                  <>
                    <div className={`text-sm font-medium mb-1 ${
                      isToday(date) ? 'text-blue-600' : 'text-gray-900'
                    }`}>
                      {date.getDate()}
                    </div>
                    
                    {/* Event indicators */}
                    <div className="space-y-1">
                      {dayEvents.slice(0, 2).map((event) => (
                        <div
                          key={event.id}
                          className={`text-xs p-1 rounded truncate border ${event.color}`}
                          title={`${event.title} - ${event.time}`}
                        >
                          {event.title}
                        </div>
                      ))}
                      {dayEvents.length > 2 && (
                        <div className="text-xs text-gray-500 px-1">
                          +{dayEvents.length - 2} more
                        </div>
                      )}
                    </div>
                  </>
                )}
              </div>
            );
          })}
        </div>
      </div>

      {/* Selected Date Events */}
      {getEventsForDate(selectedDate).length > 0 && (
        <div className="border-t border-purple-100 p-6 bg-gradient-to-r from-purple-50 to-blue-50">
          <h4 className="text-lg font-semibold text-gray-900 mb-4">
            Events for {selectedDate.toLocaleDateString()}
          </h4>
          <div className="space-y-3">
            {getEventsForDate(selectedDate).map((event) => (
              <div key={event.id} className={`p-4 rounded-xl border ${event.color} hover:shadow-md transition-shadow duration-200`}>
                <div className="flex items-start justify-between">
                  <div className="flex items-start space-x-3">
                    <div className="mt-1">
                      <EventTypeIcon type={event.type} />
                    </div>
                    <div>
                      <h5 className="font-semibold">{event.title}</h5>
                      <p className="text-sm opacity-75 mb-2">{event.course}</p>
                      <div className="flex items-center space-x-4 text-sm opacity-75">
                        <div className="flex items-center space-x-1">
                          <Clock className="w-3 h-3" />
                          <span>{event.time} ({event.duration})</span>
                        </div>
                        <div className="flex items-center space-x-1">
                          <MapPin className="w-3 h-3" />
                          <span>{event.location}</span>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Quick Add Event Modal */}
      {showEventModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl max-w-md w-full p-6">
            <h3 className="text-xl font-bold text-gray-900 mb-4">Add New Event</h3>
            <div className="space-y-4">
              <input
                type="text"
                placeholder="Event title"
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
              />
              <input
                type="date"
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
              />
              <input
                type="time"
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
              />
              <select className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500">
                <option>Select course</option>
                <option>Citizenship Prep</option>
                <option>English for Citizenship</option>
                <option>Interview Prep</option>
              </select>
            </div>
            <div className="flex space-x-3 mt-6">
              <button
                onClick={() => setShowEventModal(false)}
                className="flex-1 px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors duration-200"
              >
                Cancel
              </button>
              <button
                onClick={() => setShowEventModal(false)}
                className="flex-1 px-4 py-2 bg-gradient-to-r from-purple-600 to-blue-600 text-white rounded-lg hover:from-purple-700 hover:to-blue-700 transition-all duration-200"
              >
                Add Event
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Calendar;
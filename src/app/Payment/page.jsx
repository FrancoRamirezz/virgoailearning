"use client"
import React, { useState } from 'react';
import { CheckCircle, Crown, Zap, Shield, X, Info } from 'lucide-react';

const PricingPage = () => {
  // Add these missing state variables
  const [selectedPlan, setSelectedPlan] = useState('Pro');
  const [billingPeriod, setBillingPeriod] = useState('monthly');
  const [showTooltip, setShowTooltip] = useState(null);

  const pricingPlans = [
    {
      name: 'Basic',
      monthlyPrice: 29,
      yearlyPrice: 290,
      period: 'month',
      savings: '2 months free',
      features: [
        { name: 'Access to 5 courses', tooltip: 'Includes beginner-level citizenship prep courses' },
        { name: 'Practice tests', tooltip: 'Unlimited practice tests with basic explanations' },
        { name: 'Email support', tooltip: 'Response within 24-48 hours' },
        { name: 'Progress tracking', tooltip: 'Basic analytics and completion tracking' },
        { name: 'Mobile app access', tooltip: 'Study on-the-go with our mobile app' }
      ],
      recommended: false,
      icon: Shield,
      color: 'from-gray-600 to-gray-700',
      bgColor: 'bg-gray-100'
    },
    {
      name: 'Pro',
      monthlyPrice: 49,
      yearlyPrice: 490,
      period: 'month',
      savings: '2 months free',
      features: [
        { name: 'Access to all courses', tooltip: 'All current and future courses included' },
        { name: 'Live tutoring sessions', tooltip: '2 sessions per month with certified instructors' },
        { name: 'Priority support', tooltip: 'Response within 2-4 hours during business days' },
        { name: 'Advanced analytics', tooltip: 'Detailed performance insights and recommendations' },
        { name: 'Mock interviews', tooltip: 'AI-powered practice interviews with feedback' },
        { name: 'Personalized study plan', tooltip: 'Custom learning path based on your progress' }
      ],
      recommended: true,
      icon: Crown,
      color: 'from-purple-600 to-blue-600',
      bgColor: 'bg-purple-100'
    },
    {
      name: 'Premium',
      monthlyPrice: 89,
      yearlyPrice: 890,
      period: 'month',
      savings: '2 months free',
      features: [
        { name: 'Everything in Pro', tooltip: 'All Pro features included' },
        { name: 'One-on-one coaching', tooltip: 'Weekly 1-hour sessions with immigration expert' },
        { name: '24/7 support', tooltip: 'Round-the-clock support via chat and phone' },
        { name: 'Custom learning path', tooltip: 'Fully personalized curriculum based on your background' },
        { name: 'Guaranteed results', tooltip: 'Pass guarantee or money back' },
        { name: 'Exam fee refund*', tooltip: 'We cover your exam fees if you don\'t pass within 6 months' }
      ],
      recommended: false,
      icon: Zap,
      color: 'from-yellow-600 to-orange-600',
      bgColor: 'bg-yellow-100'
    }
  ];

  const getPrice = (plan) => {
    return billingPeriod === 'monthly' ? plan.monthlyPrice : plan.yearlyPrice;
  };

  const handlePlanSelect = (planName) => {
    setSelectedPlan(planName);
  };

  const handleTooltip = (featureIndex, planIndex) => {
    const tooltipId = `${planIndex}-${featureIndex}`;
    setShowTooltip(showTooltip === tooltipId ? null : tooltipId);
  };

  return (
    <div className="min-h-screen bg-gradient-to-b from-purple-100 to-blue-100 py-12">
      <div className="px-4 mx-auto max-w-7xl lg:px-8">
        <div className="text-center mb-12">
          <h1 className="text-4xl font-bold text-gray-900 mb-4">
            Choose Your Learning Plan
          </h1>
          <p className="text-xl text-gray-700 max-w-3xl mx-auto mb-8">
            Select the perfect plan for your citizenship test preparation journey
          </p>
          
          {/* Billing Toggle */}
          <div className="flex items-center justify-center gap-4 mb-8">
            <span className={`font-medium ${billingPeriod === 'monthly' ? 'text-gray-900' : 'text-gray-500'}`}>
              Monthly
            </span>
            <button
              onClick={() => setBillingPeriod(billingPeriod === 'monthly' ? 'yearly' : 'monthly')}
              className="relative w-14 h-7 bg-gray-300 rounded-full transition-colors duration-300 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:ring-offset-2"
              style={{
                backgroundColor: billingPeriod === 'yearly' ? '#8B5CF6' : '#D1D5DB'
              }}
              aria-label="Toggle billing period"
            >
              <div 
                className="absolute top-1 w-5 h-5 bg-white rounded-full transition-transform duration-300 shadow-md"
                style={{
                  transform: billingPeriod === 'yearly' ? 'translateX(28px)' : 'translateX(4px)'
                }}
              />
            </button>
            <span className={`font-medium ${billingPeriod === 'yearly' ? 'text-gray-900' : 'text-gray-500'}`}>
              Yearly
            </span>
            {billingPeriod === 'yearly' && (
              <span className="ml-2 px-2 py-1 text-xs font-semibold text-green-700 bg-green-100 rounded-full">
                Save 20%
              </span>
            )}
          </div>
        </div>
        
        <div className="grid gap-8 md:grid-cols-3 max-w-6xl mx-auto">
          {pricingPlans.map((plan, planIndex) => (
            <div 
              key={planIndex} 
              className={`bg-white/90 backdrop-blur-sm rounded-2xl shadow-lg hover:shadow-2xl transition-all duration-300 overflow-hidden relative border cursor-pointer transform hover:scale-105 ${
                plan.recommended 
                  ? 'ring-4 ring-purple-500 ring-opacity-50 scale-105 border-purple-200' 
                  : selectedPlan === plan.name 
                    ? 'ring-2 ring-purple-300 border-purple-200' 
                    : 'border-purple-100 hover:border-purple-200'
              }`}
              onClick={() => handlePlanSelect(plan.name)}
              role="button"
              tabIndex={0}
              onKeyDown={(e) => {
                if (e.key === 'Enter' || e.key === ' ') {
                  handlePlanSelect(plan.name);
                }
              }}
              aria-label={`Select ${plan.name} plan`}
            >
              {plan.recommended && (
                <div className={`absolute top-0 left-0 right-0 bg-gradient-to-r ${plan.color} text-white text-center py-2 text-sm font-semibold`}>
                  Most Popular
                </div>
              )}
              
              {selectedPlan === plan.name && !plan.recommended && (
                <div className="absolute top-0 left-0 right-0 bg-gradient-to-r from-green-500 to-green-600 text-white text-center py-2 text-sm font-semibold">
                  Selected Plan
                </div>
              )}
              
              <div className={`p-8 ${(plan.recommended || selectedPlan === plan.name) ? 'pt-12' : ''}`}>
                <div className="text-center mb-8">
                  <div className={`w-16 h-16 bg-gradient-to-r ${plan.color} rounded-xl flex items-center justify-center mx-auto mb-4 transition-transform duration-300 hover:scale-110`}>
                    <plan.icon className="w-8 h-8 text-white" />
                  </div>
                  <h3 className="text-2xl font-bold text-gray-900 mb-2">{plan.name}</h3>
                  <div className="mb-4">
                    <span className="text-4xl font-bold text-gray-900">
                      ${getPrice(plan)}
                    </span>
                    <span className="text-gray-600">
                      /{billingPeriod === 'monthly' ? 'month' : 'year'}
                    </span>
                    {billingPeriod === 'yearly' && (
                      <div className="text-sm text-green-600 font-medium mt-1">
                        {plan.savings}
                      </div>
                    )}
                  </div>
                </div>
                
                <ul className="space-y-4 mb-8">
                  {plan.features.map((feature, featureIndex) => (
                    <li key={featureIndex} className="flex items-start gap-3">
                      <CheckCircle className="w-5 h-5 text-green-500 flex-shrink-0 mt-0.5" />
                      <div className="flex-1">
                        <span className="text-gray-700">{feature.name}</span>
                        {feature.tooltip && (
                          <button
                            onClick={(e) => {
                              e.stopPropagation();
                              handleTooltip(featureIndex, planIndex);
                            }}
                            className="ml-2 text-gray-400 hover:text-gray-600 transition-colors"
                            aria-label={`More info about ${feature.name}`}
                          >
                            <Info className="w-4 h-4" />
                          </button>
                        )}
                        {showTooltip === `${planIndex}-${featureIndex}` && (
                          <div className="mt-2 p-3 bg-gray-100 rounded-lg text-sm text-gray-600 relative">
                            <button
                              onClick={(e) => {
                                e.stopPropagation();
                                setShowTooltip(null);
                              }}
                              className="absolute top-1 right-1 text-gray-400 hover:text-gray-600"
                              aria-label="Close tooltip"
                            >
                              <X className="w-3 h-3" />
                            </button>
                            {feature.tooltip}
                          </div>
                        )}
                      </div>
                    </li>
                  ))}
                </ul>
                
                <button 
                  className={`w-full py-3 px-6 rounded-lg font-semibold transition-all duration-300 focus:outline-none focus:ring-2 focus:ring-offset-2 ${
                    plan.recommended 
                      ? `bg-gradient-to-r ${plan.color} text-white hover:opacity-90 shadow-lg hover:shadow-xl focus:ring-purple-500` 
                      : selectedPlan === plan.name
                        ? 'bg-gradient-to-r from-green-500 to-green-600 text-white hover:from-green-600 hover:to-green-700 focus:ring-green-500'
                        : 'bg-gray-100 text-gray-900 hover:bg-gray-200 focus:ring-gray-500'
                  }`}
                  onClick={(e) => {
                    e.stopPropagation();
                    // Handle subscription logic here
                    alert(`Starting ${plan.name} plan subscription!`);
                  }}
                >
                  {selectedPlan === plan.name ? 'Start Selected Plan' : 'Get Started'}
                </button>
              </div>
            </div>
          ))}
        </div>
        
        {/* Additional Information */}
        <div className="text-center mt-12">
          <div className="bg-white/80 backdrop-blur-sm rounded-2xl p-8 border border-purple-100 max-w-4xl mx-auto">
            <h3 className="text-2xl font-bold text-gray-900 mb-4">Frequently Asked Questions</h3>
            
            <div className="grid gap-6 md:grid-cols-2 text-left">
              <div>
                <h4 className="font-semibold text-gray-900 mb-2">Can I change plans anytime?</h4>
                <p className="text-gray-600 text-sm">Yes! You can upgrade or downgrade your plan at any time. Changes take effect immediately.</p>
              </div>
              
              <div>
                <h4 className="font-semibold text-gray-900 mb-2">Is there a money-back guarantee?</h4>
                <p className="text-gray-600 text-sm">Absolutely! All plans include a 30-day money-back guarantee, no questions asked.</p>
              </div>
              
              <div>
                <h4 className="font-semibold text-gray-900 mb-2">What payment methods do you accept?</h4>
                <p className="text-gray-600 text-sm">We accept all major credit cards, PayPal, and bank transfers for yearly plans.</p>
              </div>
              
              <div>
                <h4 className="font-semibold text-gray-900 mb-2">Do you offer student discounts?</h4>
                <p className="text-gray-600 text-sm">Yes! Students get 20% off all plans with valid student ID verification.</p>
              </div>
            </div>
          </div>
          
          <div className="mt-8">
            <p className="text-gray-700 mb-4 font-medium">
              All plans include a 30-day money-back guarantee
            </p>
            <div className="text-sm text-gray-600 space-y-1">
              <p>* Premium plan includes exam fee refund if you don't pass within 6 months</p>
              <p>** Pricing in USD. Taxes may apply based on your location</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default PricingPage;
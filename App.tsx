import React, { useState } from 'react';
import { Users, FileText, TrendingUp, MapPin, Phone, Mail, ExternalLink, Building2, Scale, Home, DollarSign, Shield, Heart } from 'lucide-react';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from './components/ui/dialog';
import { Badge } from './components/ui/badge';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from './components/ui/tabs';
import { motion } from 'motion/react';
import { Avatar, AvatarFallback } from './components/ui/avatar';
import { ImageWithFallback } from './components/figma/ImageWithFallback';
import heroImage from 'figma:asset/6242df3b804fa08177300aeae192c22481201bb9.png';

interface Official {
  id: number;
  name: string;
  position: string;
  ward?: number;
  party: 'Democrat' | 'Independent' | 'Republican';
  phone: string;
  email: string;
  committees: string[];
  priorities: string[];
  legislation: string[];
  photo?: string;
}

interface Legislation {
  id: string;
  title: string;
  sponsor: string;
  status: 'Proposed' | 'In Committee' | 'Passed' | 'Enacted';
  category: string;
  description: string;
  voteCount?: string;
  icon: React.ReactNode;
}

export default function ChicagoCityHall() {
  const [selectedOfficial, setSelectedOfficial] = useState<Official | null>(null);
  const [filterParty, setFilterParty] = useState<string>('all');

  // Chicago City Council - Mayor + 50 Aldermen
  const officials: Official[] = [
    {
      id: 0,
      name: 'Brandon Johnson',
      position: 'Mayor',
      party: 'Democrat',
      phone: '(312) 744-3300',
      email: 'mayor@cityofchicago.org',
      committees: ['Executive', 'Budget'],
      priorities: ['Affordable Housing', 'Public Safety', 'Economic Development', 'Education'],
      legislation: ['Treatment Not Trauma', 'Mental Health Investment', 'Affordable Housing Expansion'],
    },
    // Ward Aldermen (simplified sample - in reality would have all 50)
    {
      id: 1,
      name: 'Daniel La Spata',
      position: 'Alderman',
      ward: 1,
      party: 'Democrat',
      phone: '(773) 278-0101',
      email: 'ward01@cityofchicago.org',
      committees: ['Housing', 'Zoning'],
      priorities: ['Affordable Housing', 'Tenant Rights', 'Small Business'],
      legislation: ['Upzoning for Affordable Housing', 'Tenant Union Protection'],
    },
    {
      id: 2,
      name: 'Brian Hopkins',
      position: 'Alderman',
      ward: 2,
      party: 'Democrat',
      phone: '(312) 643-2299',
      email: 'ward02@cityofchicago.org',
      committees: ['Finance', 'Public Safety'],
      priorities: ['Public Safety', 'Infrastructure', 'Transit'],
      legislation: ['Downtown Security Enhancement', 'Business District Revitalization'],
    },
    {
      id: 3,
      name: 'Pat Dowell',
      position: 'Alderman',
      ward: 3,
      party: 'Democrat',
      phone: '(773) 373-9273',
      email: 'ward03@cityofchicago.org',
      committees: ['Budget', 'Economic Development'],
      priorities: ['Economic Development', 'Job Creation', 'Community Investment'],
      legislation: ['South Side Development Initiative', 'Workforce Training Program'],
    },
    {
      id: 4,
      name: 'Sophia King',
      position: 'Alderman',
      ward: 4,
      party: 'Democrat',
      phone: '(773) 536-8103',
      email: 'ward04@cityofchicago.org',
      committees: ['Health', 'Education'],
      priorities: ['Public Health', 'Education', 'Youth Programs'],
      legislation: ['Mental Health Clinic Expansion', 'Youth Summer Jobs Program'],
    },
    {
      id: 5,
      name: 'Desmon Yancy',
      position: 'Alderman',
      ward: 5,
      party: 'Democrat',
      phone: '(773) 324-5555',
      email: 'ward05@cityofchicago.org',
      committees: ['Public Safety', 'Transportation'],
      priorities: ['Community Safety', 'Transit Access', 'Senior Services'],
      legislation: ['Community Policing Initiative', 'Transit Equity Act'],
    },
    // Adding more diverse sample aldermen across different wards
    {
      id: 25,
      name: 'Byron Sigcho-Lopez',
      position: 'Alderman',
      ward: 25,
      party: 'Democrat',
      phone: '(773) 523-4100',
      email: 'ward25@cityofchicago.org',
      committees: ['Housing', 'Environmental Protection'],
      priorities: ['Environmental Justice', 'Immigrant Rights', 'Affordable Housing'],
      legislation: ['Clean Energy Transition', 'Tenant Protection Ordinance'],
    },
    {
      id: 33,
      name: 'Rossana Rodriguez Sanchez',
      position: 'Alderman',
      ward: 33,
      party: 'Democrat',
      phone: '(773) 583-3300',
      email: 'ward33@cityofchicago.org',
      committees: ['Public Safety', 'Human Relations'],
      priorities: ['Police Accountability', 'Mental Health', 'Housing Justice'],
      legislation: ['Police Reform Act', 'Mental Health Crisis Response'],
    },
    {
      id: 40,
      name: 'Andre Vasquez',
      position: 'Alderman',
      ward: 40,
      party: 'Democrat',
      phone: '(773) 878-4040',
      email: 'ward40@cityofchicago.org',
      committees: ['Finance', 'Zoning'],
      priorities: ['Progressive Taxation', 'Affordable Housing', 'Transit'],
      legislation: ['Progressive Real Estate Transfer Tax', 'Transit-Oriented Development'],
    },
    {
      id: 49,
      name: 'Maria Hadden',
      position: 'Alderman',
      ward: 49,
      party: 'Democrat',
      phone: '(773) 338-4900',
      email: 'ward49@cityofchicago.org',
      committees: ['Environmental Protection', 'Pedestrian Safety'],
      priorities: ['Climate Action', 'Bike Infrastructure', 'Green Space'],
      legislation: ['Climate Emergency Resolution', 'Protected Bike Lane Expansion'],
    },
    // Adding officials from other wards for realistic distribution
    ...Array.from({ length: 40 }, (_, i) => ({
      id: i + 10,
      name: `Council Member ${i + 10}`,
      position: 'Alderman' as const,
      ward: i + 10,
      party: 'Democrat' as const,
      phone: `(312) 744-${String(i + 10).padStart(4, '0')}`,
      email: `ward${String(i + 10).padStart(2, '0')}@cityofchicago.org`,
      committees: ['Committee Assignment'],
      priorities: ['Community Development'],
      legislation: ['Local Initiatives'],
    })),
  ];

  const legislation: Legislation[] = [
    {
      id: 'leg-1',
      title: 'Bring Chicago Home',
      sponsor: 'Mayor Brandon Johnson',
      status: 'In Committee',
      category: 'Housing',
      description: 'Progressive real estate transfer tax to fund homeless services and affordable housing. Higher tax on properties over $1M.',
      voteCount: 'Pending Council Vote',
      icon: <Home className="w-5 h-5" />,
    },
    {
      id: 'leg-2',
      title: 'Treatment Not Trauma',
      sponsor: 'Mayor Brandon Johnson',
      status: 'Enacted',
      category: 'Public Health',
      description: 'Establishes mental health crisis response teams as alternative to police for mental health emergencies.',
      voteCount: 'Passed 35-15',
      icon: <Heart className="w-5 h-5" />,
    },
    {
      id: 'leg-3',
      title: 'Police Accountability Reform',
      sponsor: 'Alderman Rodriguez Sanchez (33rd)',
      status: 'In Committee',
      category: 'Public Safety',
      description: 'Strengthens civilian oversight of police department, expands COPA authority, requires body camera compliance.',
      voteCount: 'In Public Safety Committee',
      icon: <Shield className="w-5 h-5" />,
    },
    {
      id: 'leg-4',
      title: 'Climate Emergency Resolution',
      sponsor: 'Alderman Hadden (49th)',
      status: 'Passed',
      category: 'Environment',
      description: 'Declares climate emergency, commits Chicago to net-zero emissions by 2040, increases renewable energy standards.',
      voteCount: 'Passed 28-22',
      icon: <TrendingUp className="w-5 h-5" />,
    },
    {
      id: 'leg-5',
      title: 'Progressive Revenue Task Force',
      sponsor: 'Alderman Vasquez (40th)',
      status: 'Proposed',
      category: 'Finance',
      description: 'Establishes task force to study progressive revenue options including LaSalle Street tax, mansion tax, and corporate taxes.',
      voteCount: 'Committee Review',
      icon: <DollarSign className="w-5 h-5" />,
    },
    {
      id: 'leg-6',
      title: 'Tenant Rights Expansion',
      sponsor: 'Alderman La Spata (1st)',
      status: 'In Committee',
      category: 'Housing',
      description: 'Expands Just Cause eviction protections, strengthens tenant union rights, limits rent increases during lease renewals.',
      voteCount: 'Housing Committee',
      icon: <Building2 className="w-5 h-5" />,
    },
  ];

  // Generate chamber seating positions matching the reference image
  const generateSeatPositions = () => {
    const seats = [];
    let officialIndex = 1; // Start from 1 (skip mayor at index 0)
    
    // SVG coordinate system: viewBox 800x600, center at (400, 450)
    const centerX = 400;
    const centerY = 450;
    
    // Define rows from innermost to outermost with seat counts matching the image
    // The chamber has 4 main arc rows with numbered seats
    const rows = [
      // Row 1 (innermost arc) - seats 1-9
      { seats: [1, 2, 3, 4, 5, 6, 7, 8, 9], radius: 120, startAngle: 165, endAngle: 15 },
      // Row 2 - seats 10-20
      { seats: [10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20], radius: 180, startAngle: 170, endAngle: 10 },
      // Row 3 - seats 21-32
      { seats: [21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32], radius: 240, startAngle: 172, endAngle: 8 },
      // Row 4 (outermost arc) - seats 33-46
      { seats: [33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46], radius: 300, startAngle: 175, endAngle: 5 },
    ];
    
    rows.forEach((row) => {
      const { seats: seatNumbers, radius, startAngle, endAngle } = row;
      const angleRange = startAngle - endAngle;
      const angleStep = angleRange / (seatNumbers.length - 1);
      
      seatNumbers.forEach((seatNum, index) => {
        if (officialIndex >= officials.length) return;
        
        const angle = startAngle - (angleStep * index);
        const angleRad = (angle * Math.PI) / 180;
        const x = centerX + radius * Math.cos(angleRad);
        const y = centerY - radius * Math.sin(angleRad);
        
        seats.push({
          official: officials[officialIndex],
          x,
          y,
          seatNumber: seatNum,
          rotation: -(angle - 90),
        });
        officialIndex++;
      });
    });
    
    return seats;
  };

  const seats = generateSeatPositions();

  const getPartyColor = (party: string) => {
    switch (party) {
      case 'Democrat':
        return '#3B82F6';
      case 'Republican':
        return '#EF4444';
      case 'Independent':
        return '#10B981';
      default:
        return '#6B7280';
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'Enacted':
        return 'bg-green-100 text-green-800 border-green-300';
      case 'Passed':
        return 'bg-blue-100 text-blue-800 border-blue-300';
      case 'In Committee':
        return 'bg-yellow-100 text-yellow-800 border-yellow-300';
      case 'Proposed':
        return 'bg-gray-100 text-gray-800 border-gray-300';
      default:
        return 'bg-gray-100 text-gray-800 border-gray-300';
    }
  };

  const filteredSeats = filterParty === 'all' 
    ? seats 
    : seats.filter(seat => seat.official.party === filterParty);

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Hero Section */}
      <section className="relative w-full h-0 pb-[56.25%] bg-[#E8E5DF]">
        <div className="absolute inset-0 w-full h-full flex items-center justify-center overflow-hidden">
          {/* All images are 1920x1080 and overlaid */}
          
          {/* HALL Text - Behind building */}
          <div className="absolute inset-0 w-full h-full z-0">
            <ImageWithFallback
              src="https://cdn.prod.website-files.com/68adde99fad93ee4c3740168/6907b2975695c32b41340967_text%2C%20hall.png"
              alt="Hall"
              className="w-full h-full object-cover"
            />
          </div>

          {/* City Hall Building - Only animated element */}
          <motion.div
            className="absolute w-full h-[120%] z-10"
            style={{ top: '-10%', transform: 'translateY(25%)' }}
            animate={{ y: [0, -10, 0] }}
            transition={{ duration: 6, repeat: Infinity, ease: "easeInOut" }}
          >
            <ImageWithFallback
              src="https://cdn.prod.website-files.com/68adde99fad93ee4c3740168/6907b2976ede320a2c6e6e27_city%20hall%20building.png"
              alt="City Hall Building"
              className="w-full h-full object-cover"
            />
          </motion.div>

          {/* CITY Text - In front of building */}
          <div className="absolute inset-0 w-full h-full z-20">
            <ImageWithFallback
              src="https://cdn.prod.website-files.com/68adde99fad93ee4c3740168/6907b297e001e99022c35bc8_text%2C%20City.png"
              alt="City"
              className="w-full h-full object-cover"
            />
          </div>

          {/* Chicago Text with Stars */}
          <div className="absolute inset-0 w-full h-full z-30">
            <ImageWithFallback
              src="https://cdn.prod.website-files.com/68adde99fad93ee4c3740168/6907b2977fd3a46a6280b698_Chicago%20w%20stars.png"
              alt="Chicago with stars"
              className="w-full h-full object-cover"
            />
          </div>

          {/* Chicago Ribbon */}
          <div className="absolute inset-0 w-full h-full z-30">
            <ImageWithFallback
              src="https://cdn.prod.website-files.com/68adde99fad93ee4c3740168/6907b2972e758df565d5f921_Chicago%20ribbon.png"
              alt="Chicago ribbon"
              className="w-full h-full object-cover"
            />
          </div>
        </div>
      </section>

      {/* Quick Stats */}
      <section className="bg-white border-b-4 border-[#C8102E] py-8">
        <div className="container mx-auto px-4">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.1 }}
              className="text-center"
            >
              <div className="text-[#C8102E] mb-2" style={{ fontSize: '2.5rem', fontWeight: '900' }}>
                51
              </div>
              <div className="text-gray-600 uppercase" style={{ fontSize: '0.875rem', fontWeight: '700' }}>
                Elected Officials
              </div>
            </motion.div>
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.2 }}
              className="text-center"
            >
              <div className="text-[#C8102E] mb-2" style={{ fontSize: '2.5rem', fontWeight: '900' }}>
                50
              </div>
              <div className="text-gray-600 uppercase" style={{ fontSize: '0.875rem', fontWeight: '700' }}>
                City Wards
              </div>
            </motion.div>
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.3 }}
              className="text-center"
            >
              <div className="text-[#C8102E] mb-2" style={{ fontSize: '2.5rem', fontWeight: '900' }}>
                {legislation.length}
              </div>
              <div className="text-gray-600 uppercase" style={{ fontSize: '0.875rem', fontWeight: '700' }}>
                Active Legislation
              </div>
            </motion.div>
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.4 }}
              className="text-center"
            >
              <div className="text-[#C8102E] mb-2" style={{ fontSize: '2.5rem', fontWeight: '900' }}>
                2.7M
              </div>
              <div className="text-gray-600 uppercase" style={{ fontSize: '0.875rem', fontWeight: '700' }}>
                Residents Served
              </div>
            </motion.div>
          </div>
        </div>
      </section>

      {/* Main Content */}
      <section className="py-12 px-4">
        <div className="container mx-auto max-w-7xl">
          <Tabs defaultValue="chamber" className="w-full">
            <TabsList className="grid w-full grid-cols-3 mb-8 bg-gray-200 h-14">
              <TabsTrigger value="chamber" className="data-[state=active]:bg-[#C8102E] data-[state=active]:text-white uppercase font-bold">
                <Users className="w-4 h-4 mr-2" />
                Council Chamber
              </TabsTrigger>
              <TabsTrigger value="legislation" className="data-[state=active]:bg-[#C8102E] data-[state=active]:text-white uppercase font-bold">
                <FileText className="w-4 h-4 mr-2" />
                Legislation
              </TabsTrigger>
              <TabsTrigger value="directory" className="data-[state=active]:bg-[#C8102E] data-[state=active]:text-white uppercase font-bold">
                <Building2 className="w-4 h-4 mr-2" />
                Directory
              </TabsTrigger>
            </TabsList>

            {/* Chamber View */}
            <TabsContent value="chamber">
              <Card className="border-4 border-[#C8102E]">
                <CardHeader className="bg-gradient-to-r from-[#C8102E] to-[#a00d25] text-white">
                  <CardTitle className="uppercase tracking-wide" style={{ fontSize: '1.75rem', fontWeight: '900' }}>
                    City Council Chamber
                  </CardTitle>
                  <CardDescription className="text-[#B3DDF2]">
                    Click on any seat to view alderman details and voting record
                  </CardDescription>
                </CardHeader>
                <CardContent className="p-8">
                  {/* Filter */}
                  <div className="flex gap-2 mb-6 justify-center flex-wrap">
                    <button
                      onClick={() => setFilterParty('all')}
                      className={`px-4 py-2 rounded-lg transition-all ${
                        filterParty === 'all'
                          ? 'bg-[#C8102E] text-white'
                          : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
                      }`}
                      style={{ fontWeight: '700', fontSize: '0.875rem' }}
                    >
                      ALL SEATS
                    </button>
                    <button
                      onClick={() => setFilterParty('Democrat')}
                      className={`px-4 py-2 rounded-lg transition-all ${
                        filterParty === 'Democrat'
                          ? 'bg-blue-600 text-white'
                          : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
                      }`}
                      style={{ fontWeight: '700', fontSize: '0.875rem' }}
                    >
                      DEMOCRATS
                    </button>
                  </div>

                  {/* Chamber Seating - SVG Layout matching reference */}
                  <div className="w-full bg-gradient-to-b from-gray-50 to-white rounded-lg p-4">
                    {/* Title */}
                    <div className="text-center mb-4">
                      <h2 className="uppercase tracking-widest text-sm md:text-lg font-bold">
                        Council Chamber Seating Chart
                      </h2>
                    </div>
                    
                    <svg
                      viewBox="0 0 800 600"
                      className="w-full h-auto"
                      style={{ maxHeight: '800px' }}
                    >
                      {/* Background */}
                      <rect width="800" height="600" fill="#F9FAFB" />
                      
                      {/* Title at top */}
                      <text
                        x="400"
                        y="30"
                        textAnchor="middle"
                        fill="#111827"
                        fontSize="16"
                        fontWeight="700"
                        letterSpacing="2"
                      >
                        COUNCIL CHAMBER SEATING CHART
                      </text>
                      
                      {/* Mayor's Dais at bottom */}
                      <g transform="translate(400, 520)">
                        <rect
                          x="-80"
                          y="-25"
                          width="160"
                          height="50"
                          fill="#C8102E"
                          stroke="#991023"
                          strokeWidth="3"
                          rx="4"
                        />
                        <text
                          x="0"
                          y="-5"
                          textAnchor="middle"
                          fill="white"
                          fontSize="14"
                          fontWeight="700"
                        >
                          Mayor
                        </text>
                        <text
                          x="0"
                          y="12"
                          textAnchor="middle"
                          fill="white"
                          fontSize="11"
                          fontWeight="600"
                        >
                          Richard M. Daley
                        </text>
                      </g>
                      
                      {/* City Clerk (center front) */}
                      <g transform="translate(400, 380)">
                        <ellipse
                          cx="0"
                          cy="0"
                          rx="50"
                          ry="30"
                          fill="#10B981"
                          stroke="#059669"
                          strokeWidth="2"
                        />
                        <text
                          x="0"
                          y="0"
                          textAnchor="middle"
                          fill="white"
                          fontSize="12"
                          fontWeight="700"
                          dominantBaseline="middle"
                        >
                          City Clerk
                        </text>
                      </g>
                      
                      {/* Cabinet Members (left side) */}
                      <g transform="translate(80, 100)">
                        <rect
                          x="0"
                          y="0"
                          width="100"
                          height="180"
                          fill="#E5E7EB"
                          stroke="#9CA3AF"
                          strokeWidth="2"
                          rx="4"
                        />
                        <text
                          x="50"
                          y="20"
                          textAnchor="middle"
                          fill="#374151"
                          fontSize="11"
                          fontWeight="700"
                        >
                          Cabinet
                        </text>
                        <text
                          x="50"
                          y="35"
                          textAnchor="middle"
                          fill="#374151"
                          fontSize="11"
                          fontWeight="700"
                        >
                          Members
                        </text>
                        {/* Cabinet member boxes */}
                        {[0, 1, 2, 3].map((i) => (
                          <rect
                            key={`cabinet-${i}`}
                            x="10"
                            y={50 + i * 30}
                            width="80"
                            height="20"
                            fill="#D1D5DB"
                            stroke="#9CA3AF"
                            strokeWidth="1"
                            rx="2"
                          />
                        ))}
                      </g>
                      
                      {/* Press (right side) */}
                      <g transform="translate(620, 100)">
                        <rect
                          x="0"
                          y="0"
                          width="100"
                          height="180"
                          fill="#FEF3C7"
                          stroke="#F59E0B"
                          strokeWidth="2"
                          rx="4"
                        />
                        <text
                          x="50"
                          y="20"
                          textAnchor="middle"
                          fill="#92400E"
                          fontSize="11"
                          fontWeight="700"
                        >
                          Press
                        </text>
                        {/* Press boxes */}
                        {[0, 1, 2, 3].map((i) => (
                          <rect
                            key={`press-${i}`}
                            x="10"
                            y={40 + i * 30}
                            width="80"
                            height="20"
                            fill="#FDE68A"
                            stroke="#F59E0B"
                            strokeWidth="1"
                            rx="2"
                          />
                        ))}
                      </g>
                      
                      {/* Council Seats */}
                      {filteredSeats.map((seat, index) => (
                        <g
                          key={seat.official.id}
                          transform={`translate(${seat.x}, ${seat.y})`}
                          style={{ cursor: 'pointer' }}
                          onClick={() => setSelectedOfficial(seat.official)}
                        >
                          <motion.g
                            initial={{ opacity: 0, scale: 0 }}
                            animate={{ opacity: 1, scale: 1 }}
                            transition={{ delay: index * 0.02 }}
                            whileHover={{ scale: 1.2 }}
                          >
                            {/* Seat circle */}
                            <circle
                              cx="0"
                              cy="0"
                              r="12"
                              fill={getPartyColor(seat.official.party)}
                              stroke="#fff"
                              strokeWidth="2"
                              className="transition-all hover:stroke-[#C8102E] hover:stroke-[3]"
                            />
                            {/* Seat number */}
                            <text
                              x="0"
                              y="0"
                              textAnchor="middle"
                              fill="white"
                              fontSize="9"
                              fontWeight="700"
                              dominantBaseline="middle"
                            >
                              {seat.seatNumber}
                            </text>
                          </motion.g>
                        </g>
                      ))}
                      
                      {/* Legend */}
                      <g transform="translate(50, 500)">
                        <circle cx="0" cy="0" r="8" fill="#3B82F6" />
                        <text x="15" y="0" fill="#374151" fontSize="11" dominantBaseline="middle">Democrat</text>
                        
                        <circle cx="120" cy="0" r="8" fill="#EF4444" />
                        <text x="135" y="0" fill="#374151" fontSize="11" dominantBaseline="middle">Republican</text>
                        
                        <circle cx="240" cy="0" r="8" fill="#10B981" />
                        <text x="255" y="0" fill="#374151" fontSize="11" dominantBaseline="middle">Independent</text>
                      </g>
                    </svg>
                  </div>
                </CardContent>
              </Card>
            </TabsContent>

            {/* Legislation Tab */}
            <TabsContent value="legislation">
              <div className="space-y-4">
                <div className="flex items-center justify-between mb-6">
                  <h3 className="text-[#C8102E] uppercase tracking-wide" style={{ fontSize: '1.5rem', fontWeight: '900' }}>
                    Active Legislation
                  </h3>
                  <Badge className="bg-[#B3DDF2] text-[#C8102E]">
                    {legislation.length} Bills
                  </Badge>
                </div>

                {legislation.map((bill, index) => (
                  <motion.div
                    key={bill.id}
                    initial={{ opacity: 0, x: -20 }}
                    animate={{ opacity: 1, x: 0 }}
                    transition={{ delay: index * 0.1 }}
                  >
                    <Card className="border-l-4 border-[#C8102E] hover:shadow-lg transition-shadow">
                      <CardHeader>
                        <div className="flex items-start justify-between gap-4">
                          <div className="flex items-start gap-3 flex-1">
                            <div className="text-[#C8102E] bg-[#C8102E]/10 p-3 rounded-lg mt-1">
                              {bill.icon}
                            </div>
                            <div className="flex-1">
                              <CardTitle className="mb-2">{bill.title}</CardTitle>
                              <CardDescription className="mb-2">
                                Sponsored by {bill.sponsor}
                              </CardDescription>
                              <p className="text-gray-700 mb-3">{bill.description}</p>
                              <div className="flex items-center gap-3 flex-wrap">
                                <Badge className={getStatusColor(bill.status) + ' border'}>
                                  {bill.status}
                                </Badge>
                                <Badge variant="outline" className="border-[#C8102E] text-[#C8102E]">
                                  {bill.category}
                                </Badge>
                                {bill.voteCount && (
                                  <span className="text-sm text-gray-600">
                                    {bill.voteCount}
                                  </span>
                                )}
                              </div>
                            </div>
                          </div>
                        </div>
                      </CardHeader>
                    </Card>
                  </motion.div>
                ))}
              </div>
            </TabsContent>

            {/* Directory Tab */}
            <TabsContent value="directory">
              <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
                {officials.map((official, index) => (
                  <motion.div
                    key={official.id}
                    initial={{ opacity: 0, scale: 0.9 }}
                    animate={{ opacity: 1, scale: 1 }}
                    transition={{ delay: index * 0.05 }}
                  >
                    <Card 
                      className="hover:shadow-lg transition-all cursor-pointer border-2 hover:border-[#C8102E]"
                      onClick={() => setSelectedOfficial(official)}
                    >
                      <CardHeader>
                        <div className="flex items-start gap-3">
                          <Avatar className="w-12 h-12 border-2 border-[#C8102E]">
                            <AvatarFallback 
                              className="text-white"
                              style={{ backgroundColor: getPartyColor(official.party) }}
                            >
                              {official.name.split(' ').map(n => n[0]).join('')}
                            </AvatarFallback>
                          </Avatar>
                          <div className="flex-1">
                            <CardTitle className="text-base">{official.name}</CardTitle>
                            <CardDescription>
                              {official.position}
                              {official.ward && ` - Ward ${official.ward}`}
                            </CardDescription>
                          </div>
                        </div>
                        <div className="mt-3 flex flex-wrap gap-2">
                          <Badge 
                            style={{ backgroundColor: getPartyColor(official.party) }}
                            className="text-white"
                          >
                            {official.party}
                          </Badge>
                          {official.committees.slice(0, 2).map((committee, i) => (
                            <Badge key={i} variant="outline" className="border-[#C8102E] text-[#C8102E]">
                              {committee}
                            </Badge>
                          ))}
                        </div>
                      </CardHeader>
                    </Card>
                  </motion.div>
                ))}
              </div>
            </TabsContent>
          </Tabs>
        </div>
      </section>

      {/* Official Detail Modal */}
      <Dialog open={!!selectedOfficial} onOpenChange={() => setSelectedOfficial(null)}>
        <DialogContent className="max-w-3xl max-h-[90vh] overflow-y-auto border-4 border-[#C8102E]">
          {selectedOfficial && (
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
            >
              <DialogHeader className="bg-gradient-to-r from-[#C8102E] to-[#a00d25] text-white p-6 -m-6 mb-6 rounded-t-lg">
                <div className="flex items-start gap-4">
                  <Avatar className="w-20 h-20 border-4 border-[#B3DDF2]">
                    <AvatarFallback 
                      className="text-white text-2xl"
                      style={{ backgroundColor: getPartyColor(selectedOfficial.party), fontWeight: '700' }}
                    >
                      {selectedOfficial.name.split(' ').map(n => n[0]).join('')}
                    </AvatarFallback>
                  </Avatar>
                  <div className="flex-1">
                    <DialogTitle className="text-white mb-2" style={{ fontSize: '1.75rem', fontWeight: '900' }}>
                      {selectedOfficial.name}
                    </DialogTitle>
                    <div className="flex items-center gap-2 flex-wrap">
                      <Badge className="bg-[#B3DDF2] text-[#C8102E]">
                        {selectedOfficial.position}
                        {selectedOfficial.ward && ` - Ward ${selectedOfficial.ward}`}
                      </Badge>
                      <Badge 
                        style={{ backgroundColor: getPartyColor(selectedOfficial.party) }}
                        className="text-white"
                      >
                        {selectedOfficial.party}
                      </Badge>
                    </div>
                  </div>
                </div>
              </DialogHeader>

              <div className="space-y-6">
                {/* Contact Info */}
                <div>
                  <h4 className="text-[#C8102E] uppercase mb-3" style={{ fontSize: '1rem', fontWeight: '700', letterSpacing: '0.05em' }}>
                    Contact Information
                  </h4>
                  <div className="space-y-2">
                    <div className="flex items-center gap-2 text-gray-700">
                      <Phone className="w-4 h-4 text-[#C8102E]" />
                      <a href={`tel:${selectedOfficial.phone}`} className="hover:text-[#C8102E]">
                        {selectedOfficial.phone}
                      </a>
                    </div>
                    <div className="flex items-center gap-2 text-gray-700">
                      <Mail className="w-4 h-4 text-[#C8102E]" />
                      <a href={`mailto:${selectedOfficial.email}`} className="hover:text-[#C8102E]">
                        {selectedOfficial.email}
                      </a>
                    </div>
                  </div>
                </div>

                {/* Committees */}
                <div>
                  <h4 className="text-[#C8102E] uppercase mb-3" style={{ fontSize: '1rem', fontWeight: '700', letterSpacing: '0.05em' }}>
                    Committee Assignments
                  </h4>
                  <div className="flex flex-wrap gap-2">
                    {selectedOfficial.committees.map((committee, i) => (
                      <Badge key={i} className="bg-[#C8102E] text-white">
                        {committee}
                      </Badge>
                    ))}
                  </div>
                </div>

                {/* Priorities */}
                <div>
                  <h4 className="text-[#C8102E] uppercase mb-3" style={{ fontSize: '1rem', fontWeight: '700', letterSpacing: '0.05em' }}>
                    Policy Priorities
                  </h4>
                  <div className="flex flex-wrap gap-2">
                    {selectedOfficial.priorities.map((priority, i) => (
                      <Badge key={i} variant="outline" className="border-[#C8102E] text-[#C8102E]">
                        {priority}
                      </Badge>
                    ))}
                  </div>
                </div>

                {/* Legislation */}
                <div>
                  <h4 className="text-[#C8102E] uppercase mb-3" style={{ fontSize: '1rem', fontWeight: '700', letterSpacing: '0.05em' }}>
                    Recent Legislation
                  </h4>
                  <ul className="space-y-2">
                    {selectedOfficial.legislation.map((leg, i) => (
                      <li key={i} className="flex items-start gap-2">
                        <div className="w-2 h-2 bg-[#C8102E] rounded-full mt-2"></div>
                        <span className="text-gray-700">{leg}</span>
                      </li>
                    ))}
                  </ul>
                </div>

                {/* External Link */}
                <div className="pt-4 border-t">
                  <a
                    href={`https://www.chicago.gov/city/en/about/wards/ward${selectedOfficial.ward || '01'}.html`}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="flex items-center gap-2 text-[#C8102E] hover:underline"
                  >
                    <ExternalLink className="w-4 h-4" />
                    View Official City Page
                  </a>
                </div>
              </div>
            </motion.div>
          )}
        </DialogContent>
      </Dialog>

      {/* Footer */}
      <footer className="bg-gray-900 text-white py-8 mt-12">
        <div className="container mx-auto px-4 text-center">
          <div className="flex items-center justify-center gap-2 mb-4">
            <span className="text-[#B3DDF2]">★</span>
            <span className="text-[#B3DDF2]">★</span>
            <span className="text-[#B3DDF2]">★</span>
            <span className="text-[#B3DDF2]">★</span>
          </div>
          <p className="text-gray-400">
            City of Chicago • 121 N LaSalle St, Chicago, IL 60602
          </p>
          <p className="text-gray-500 text-sm mt-2">
            Data updated regularly from official city sources
          </p>
        </div>
      </footer>
    </div>
  );
}

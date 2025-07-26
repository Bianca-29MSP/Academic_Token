import React, { useEffect, useRef } from 'react';
import * as d3 from 'd3';
import { SemesterWithSubjects } from '../../types/curriculum.types';

interface PrerequisiteFlowProps {
  semesters: SemesterWithSubjects[];
  completedSubjects: string[];
}

interface Node {
  id: string;
  label: string;
  title: string;
  semester: number;
  x: number;
  y: number;
  completed: boolean;
}

interface Link {
  source: string;
  target: string;
}

const PrerequisiteFlow: React.FC<PrerequisiteFlowProps> = ({ 
  semesters, 
  completedSubjects 
}) => {
  const svgRef = useRef<SVGSVGElement>(null);

  useEffect(() => {
    if (!svgRef.current || semesters.length === 0) return;

    // Clear previous content
    d3.select(svgRef.current).selectAll('*').remove();

    // Create nodes and links for D3
    const nodes: Node[] = [];
    const links: Link[] = [];

    // Create nodes from subjects
    semesters.forEach((semester, semesterIndex) => {
      semester.subjects.forEach((subject, subjectIndex) => {
        nodes.push({
          id: subject.subjectId,
          label: subject.code,
          title: subject.title,
          semester: parseInt(semester.semesterNumber),
          x: semesterIndex * 200 + 100,
          y: subjectIndex * 80 + 100,
          completed: completedSubjects.includes(subject.subjectId),
        });

        // Create links for prerequisites
        if (subject.prerequisites) {
          subject.prerequisites.forEach(prereqId => {
            links.push({
              source: prereqId,
              target: subject.subjectId,
            });
          });
        }
      });
    });

    // Set up SVG dimensions
    const width = semesters.length * 200 + 200;
    const height = Math.max(...semesters.map(s => s.subjects.length)) * 80 + 200;

    const svg = d3.select(svgRef.current)
      .attr('width', width)
      .attr('height', height)
      .attr('viewBox', `0 0 ${width} ${height}`);

    // Add arrow marker definition
    svg.append('defs').append('marker')
      .attr('id', 'arrow')
      .attr('viewBox', '0 0 10 10')
      .attr('refX', 10)
      .attr('refY', 5)
      .attr('markerWidth', 5)
      .attr('markerHeight', 5)
      .attr('orient', 'auto')
      .append('path')
      .attr('d', 'M 0 0 L 10 5 L 0 10 z')
      .attr('fill', '#999');

    // Add semester labels
    const semesterGroups = svg.selectAll('.semester-group')
      .data(semesters)
      .enter()
      .append('g')
      .attr('class', 'semester-group')
      .attr('transform', (d, i) => `translate(${i * 200 + 100}, 50)`);

    semesterGroups.append('text')
      .attr('text-anchor', 'middle')
      .attr('font-weight', 'bold')
      .attr('font-size', '16px')
      .attr('fill', '#333')
      .text(d => `${d.semesterNumber}º Período`);

    // Draw links
    const link = svg.append('g')
      .selectAll('line')
      .data(links)
      .enter()
      .append('line')
      .attr('stroke', '#999')
      .attr('stroke-width', 2)
      .attr('marker-end', 'url(#arrow)');

    // Draw nodes
    const node = svg.append('g')
      .selectAll('g')
      .data(nodes)
      .enter()
      .append('g')
      .attr('transform', d => `translate(${d.x}, ${d.y})`)
      .attr('cursor', 'pointer');

    // Add rectangles for nodes
    node.append('rect')
      .attr('width', 150)
      .attr('height', 60)
      .attr('x', -75)
      .attr('y', -30)
      .attr('rx', 5)
      .attr('fill', d => d.completed ? '#4caf50' : '#f0f0f0')
      .attr('stroke', d => d.completed ? '#388e3c' : '#ccc')
      .attr('stroke-width', 2);

    // Add code label
    node.append('text')
      .attr('text-anchor', 'middle')
      .attr('dy', -10)
      .attr('font-weight', 'bold')
      .attr('font-size', '14px')
      .attr('fill', d => d.completed ? 'white' : '#333')
      .text(d => d.label);

    // Add title (truncated)
    node.append('text')
      .attr('text-anchor', 'middle')
      .attr('dy', 10)
      .attr('font-size', '12px')
      .attr('fill', d => d.completed ? 'white' : '#666')
      .text(d => d.title.length > 20 ? d.title.substring(0, 20) + '...' : d.title);

    // Update link positions
    link
      .attr('x1', d => {
        const source = nodes.find(n => n.id === d.source);
        return source ? source.x + 75 : 0;
      })
      .attr('y1', d => {
        const source = nodes.find(n => n.id === d.source);
        return source ? source.y : 0;
      })
      .attr('x2', d => {
        const target = nodes.find(n => n.id === d.target);
        return target ? target.x - 75 : 0;
      })
      .attr('y2', d => {
        const target = nodes.find(n => n.id === d.target);
        return target ? target.y : 0;
      });

    // Add hover effects
    node.on('mouseenter', function() {
      d3.select(this).select('rect')
        .transition()
        .duration(200)
        .attr('transform', 'scale(1.05)');
    })
    .on('mouseleave', function() {
      d3.select(this).select('rect')
        .transition()
        .duration(200)
        .attr('transform', 'scale(1)');
    });

  }, [semesters, completedSubjects]);

  return (
    <div className="prerequisite-flow" style={{ overflowX: 'auto', padding: '20px' }}>
      <svg ref={svgRef}></svg>
    </div>
  );
};

export default PrerequisiteFlow;

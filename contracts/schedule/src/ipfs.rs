// src/ipfs.rs
use cosmwasm_schema::cw_serde;
use cosmwasm_std::{Deps, DepsMut, StdError, StdResult};
use crate::state::IPFS_CACHE;

/// Subject content retrieved from IPFS for schedule analysis
#[cw_serde]
pub struct SubjectContent {
    pub title: String,
    pub code: String,
    pub description: String,
    pub objectives: Vec<String>,
    pub topics: Vec<TopicUnit>,
    pub competencies: Vec<String>,
    pub methodology: Vec<String>,
    pub assessment_methods: Vec<String>,
    pub workload_hours: u64,
    pub credits: u64,
    pub prerequisites_description: String,
    pub learning_outcomes: Vec<String>,
    pub knowledge_areas: Vec<String>,
    pub difficulty_level: String,
    pub language: String,
    pub institution: String,
    pub professor: String,
    pub semester: String,
    pub year: String,
    pub bibliography: Vec<String>,
    pub supplementary_materials: Vec<String>,
    pub practical_activities: Vec<String>,
    pub theoretical_hours: u64,
    pub practical_hours: u64,
    pub content_hash: String,
    
    // Schedule-specific fields
    pub class_schedule: Option<ClassScheduleDetails>,
    pub enrollment_capacity: Option<u32>,
    pub historical_demand: Option<DemandMetrics>,
    pub success_metrics: Option<SubjectSuccessMetrics>,
}

/// Topic or content unit within the subject
#[cw_serde]
pub struct TopicUnit {
    pub unit_number: u32,
    pub title: String,
    pub description: String,
    pub hours: u64,
    pub keywords: Vec<String>,
    pub subtopics: Vec<String>,
    pub learning_objectives: Vec<String>,
    pub required_knowledge: Vec<String>,
    pub complexity_rating: Option<u32>,
}

/// Detailed class schedule information
#[cw_serde]
pub struct ClassScheduleDetails {
    pub time_slots: Vec<TimeSlotDetail>,
    pub location: String,
    pub online_option: bool,
    pub hybrid_format: bool,
    pub attendance_policy: String,
    pub makeup_policy: String,
}

/// Specific time slot for classes
#[cw_serde]
pub struct TimeSlotDetail {
    pub day: String,
    pub start_time: String,
    pub end_time: String,
    pub duration_minutes: u32,
    pub frequency: String, // "weekly", "biweekly", etc.
}

/// Historical demand metrics for the subject
#[cw_serde]
pub struct DemandMetrics {
    pub enrollment_history: Vec<SemesterEnrollment>,
    pub average_demand: u32,
    pub demand_trend: DemandTrend,
    pub waitlist_frequency: u32,
}

/// Enrollment data for a specific semester
#[cw_serde]
pub struct SemesterEnrollment {
    pub semester: String,
    pub enrolled_students: u32,
    pub waitlisted_students: u32,
    pub capacity: u32,
    pub completion_rate: u32,
}

/// Trend in demand over time
#[cw_serde]
pub enum DemandTrend {
    Increasing,
    Stable,
    Decreasing,
    Seasonal,
    Unpredictable,
}

/// Success metrics for the subject
#[cw_serde]
pub struct SubjectSuccessMetrics {
    pub average_grade: u32,
    pub completion_rate: u32,
    pub student_satisfaction: u32,
    pub difficulty_rating: u32,
    pub workload_rating: u32,
    pub prerequisite_success_correlation: Vec<PrerequisiteCorrelation>,
}

/// Correlation between prerequisite performance and subject success
#[cw_serde]
pub struct PrerequisiteCorrelation {
    pub prerequisite_subject_id: String,
    pub correlation_strength: u32,
    pub minimum_recommended_grade: u32,
}

/// Fetch subject content from IPFS cache
pub fn fetch_ipfs_content(deps: Deps, ipfs_link: &str) -> StdResult<SubjectContent> {
    match IPFS_CACHE.may_load(deps.storage, ipfs_link)? {
        Some(content) => Ok(content),
        None => Err(StdError::generic_err(format!(
            "IPFS content not found in cache for link: {}. Please ensure content is cached first.", 
            ipfs_link
        ))),
    }
}

/// Cache IPFS content in the contract storage
pub fn cache_ipfs_content(
    deps: DepsMut,
    ipfs_link: &str,
    content: &SubjectContent,
) -> StdResult<()> {
    IPFS_CACHE.save(deps.storage, ipfs_link, content)
}

/// Check if IPFS content exists in cache
pub fn is_content_cached(deps: Deps, ipfs_link: &str) -> StdResult<bool> {
    Ok(IPFS_CACHE.may_load(deps.storage, ipfs_link)?.is_some())
}

/// Verify content integrity using hash
pub fn verify_content_integrity(content: &SubjectContent, expected_hash: &str) -> bool {
    calculate_content_hash(content) == expected_hash
}

/// Calculate hash of content for integrity verification
pub fn calculate_content_hash(content: &SubjectContent) -> String {
    use sha2::{Digest, Sha256};
    
    let content_str = format!(
        "{}|{}|{}|{}|{}|{}|{}|{}",
        content.title,
        content.code,
        content.description,
        content.objectives.join(","),
        content.competencies.join(","),
        content.knowledge_areas.join(","),
        content.workload_hours,
        content.credits
    );
    
    let mut hasher = Sha256::new();
    hasher.update(content_str.as_bytes());
    format!("{:x}", hasher.finalize())
}

/// Analyze subject difficulty based on content
pub fn analyze_subject_difficulty(content: &SubjectContent) -> u32 {
    let mut difficulty_score = 0u32;
    
    // Base difficulty from explicit rating
    if let Ok(explicit_difficulty) = content.difficulty_level.parse::<u32>() {
        difficulty_score += explicit_difficulty * 20;
    }
    
    // Complexity from number of topics and depth
    let topic_complexity = (content.topics.len() as u32).min(10) * 5;
    difficulty_score += topic_complexity;
    
    // Workload factor
    let workload_factor = (content.workload_hours as u32 / 15).min(20);
    difficulty_score += workload_factor;
    
    // Prerequisites complexity
    let prerequisites_count = content.prerequisites_description
        .split(',')
        .filter(|s| !s.trim().is_empty())
        .count() as u32;
    difficulty_score += prerequisites_count.min(10) * 3;
    
    // Assessment methods complexity
    let assessment_complexity = content.assessment_methods.len() as u32 * 2;
    difficulty_score += assessment_complexity.min(10);
    
    // Knowledge areas breadth
    let knowledge_breadth = content.knowledge_areas.len() as u32 * 2;
    difficulty_score += knowledge_breadth.min(15);
    
    difficulty_score.min(100)
}

/// Estimate weekly workload hours based on content
pub fn estimate_weekly_workload(content: &SubjectContent) -> u32 {
    // Base hours from credits (typically 15 hours per credit per semester)
    let credit_hours = content.credits as u32 * 15;
    let mut weekly_credit_hours = credit_hours / 16; // Assuming 16-week semester
    
    // Theoretical vs practical ratio adjustment
    let theory_ratio = if content.theoretical_hours + content.practical_hours > 0 {
        content.theoretical_hours as f64 / (content.theoretical_hours + content.practical_hours) as f64
    } else {
        0.7 // Default assumption
    };
    
    // Practical subjects typically require more study time
    let workload_multiplier = if theory_ratio > 0.8 {
        1.2 // Theory-heavy subjects need more study time
    } else if theory_ratio < 0.3 {
        0.9 // Practice-heavy subjects are more hands-on
    } else {
        1.0 // Balanced subjects
    };
    
    weekly_credit_hours = (weekly_credit_hours as f64 * workload_multiplier) as u32;
    
    // Adjust for difficulty
    let difficulty = analyze_subject_difficulty(content);
    let difficulty_multiplier = 1.0 + (difficulty as f64 / 200.0);
    weekly_credit_hours = (weekly_credit_hours as f64 * difficulty_multiplier) as u32;
    
    // Practical bounds
    weekly_credit_hours.max(2).min(25)
}

/// Analyze schedule conflicts between subjects
pub fn analyze_schedule_conflicts(
    subject1: &SubjectContent,
    subject2: &SubjectContent,
) -> Vec<String> {
    let mut conflicts = Vec::new();
    
    // Check if both subjects have schedule information
    if let (Some(sched1), Some(sched2)) = (&subject1.class_schedule, &subject2.class_schedule) {
        // Check time conflicts
        for slot1 in &sched1.time_slots {
            for slot2 in &sched2.time_slots {
                if slot1.day == slot2.day {
                    if time_slots_overlap(&slot1.start_time, &slot1.end_time, 
                                        &slot2.start_time, &slot2.end_time) {
                        conflicts.push(format!(
                            "Time conflict on {} between {}-{} and {}-{}",
                            slot1.day,
                            slot1.start_time, slot1.end_time,
                            slot2.start_time, slot2.end_time
                        ));
                    }
                }
            }
        }
        
        // Check location conflicts if same time
        if sched1.location == sched2.location && !conflicts.is_empty() {
            conflicts.push(format!("Location conflict at {}", sched1.location));
        }
    }
    
    // Check workload conflicts
    let total_workload = estimate_weekly_workload(subject1) + estimate_weekly_workload(subject2);
    if total_workload > 20 {
        conflicts.push(format!(
            "High combined workload: {} hours per week",
            total_workload
        ));
    }
    
    // Check knowledge area conflicts (too much overlap might be redundant)
    let overlap = knowledge_area_overlap(&subject1.knowledge_areas, &subject2.knowledge_areas);
    if overlap > 70 {
        conflicts.push("High knowledge area overlap - potentially redundant".to_string());
    }
    
    conflicts
}

/// Check if two time slots overlap
fn time_slots_overlap(start1: &str, end1: &str, start2: &str, end2: &str) -> bool {
    // Simple time comparison (assumes HH:MM format)
    // In a real implementation, you'd parse the time strings properly
    let start1_mins = parse_time_to_minutes(start1);
    let end1_mins = parse_time_to_minutes(end1);
    let start2_mins = parse_time_to_minutes(start2);
    let end2_mins = parse_time_to_minutes(end2);
    
    // Check for overlap
    start1_mins < end2_mins && start2_mins < end1_mins
}

/// Parse time string to minutes since midnight
fn parse_time_to_minutes(time: &str) -> u32 {
    let parts: Vec<&str> = time.split(':').collect();
    if parts.len() == 2 {
        if let (Ok(hours), Ok(minutes)) = (parts[0].parse::<u32>(), parts[1].parse::<u32>()) {
            return hours * 60 + minutes;
        }
    }
    0 // Default to midnight if parsing fails
}

/// Calculate knowledge area overlap percentage
fn knowledge_area_overlap(areas1: &[String], areas2: &[String]) -> u32 {
    if areas1.is_empty() || areas2.is_empty() {
        return 0;
    }
    
    let set1: std::collections::HashSet<_> = areas1.iter().collect();
    let set2: std::collections::HashSet<_> = areas2.iter().collect();
    
    let intersection = set1.intersection(&set2).count();
    let union = set1.union(&set2).count();
    
    if union > 0 {
        ((intersection * 100) / union) as u32
    } else {
        0
    }
}

/// Generate subject recommendations based on content analysis
pub fn generate_content_based_recommendations(
    content: &SubjectContent,
    _student_completed_subjects: &[String],
) -> Vec<String> {
    let mut recommendations = Vec::new();
    
    // Analyze difficulty vs student's track record
    let difficulty = analyze_subject_difficulty(content);
    if difficulty > 80 {
        recommendations.push("High difficulty subject - ensure strong foundation in prerequisites".to_string());
    }
    
    // Workload analysis
    let workload = estimate_weekly_workload(content);
    if workload > 15 {
        recommendations.push(format!(
            "High workload subject ({} hours/week) - consider lighter course load this semester",
            workload
        ));
    }
    
    // Prerequisites analysis
    if !content.prerequisites_description.is_empty() {
        recommendations.push("Verify all prerequisites are completed with satisfactory grades".to_string());
    }
    
    // Success metrics analysis
    if let Some(metrics) = &content.success_metrics {
        if metrics.completion_rate < 80 {
            recommendations.push("Lower than average completion rate - consider additional preparation".to_string());
        }
        
        if metrics.average_grade < 75 {
            recommendations.push("Subject has lower average grades - plan for additional study time".to_string());
        }
    }
    
    // Demand analysis
    if let Some(demand) = &content.historical_demand {
        if demand.waitlist_frequency > 50 {
            recommendations.push("High demand subject - register early and have backup options".to_string());
        }
    }
    
    recommendations
}
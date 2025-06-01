// Progress Contract - src/contract.rs
use cosmwasm_std::{
    entry_point, to_json_binary, Binary, Deps, DepsMut, Env, MessageInfo, Response, StdResult,
};
use cw2::set_contract_version;

use crate::error::ContractError;
use crate::msg::{ExecuteMsg, InstantiateMsg, QueryMsg};
use crate::state::{State, ProgressConfig, STATE, PROGRESS_CONFIG};
use crate::{execute, query};

const CONTRACT_NAME: &str = "crates.io:progress";
const CONTRACT_VERSION: &str = env!("CARGO_PKG_VERSION");

#[entry_point]
pub fn instantiate(
    deps: DepsMut,
    _env: Env,
    info: MessageInfo,
    msg: InstantiateMsg,
) -> Result<Response, ContractError> {
    set_contract_version(deps.storage, CONTRACT_NAME, CONTRACT_VERSION)?;

    let owner = msg.owner
        .and_then(|s| deps.api.addr_validate(&s).ok())
        .unwrap_or(info.sender);

    let state = State {
        owner: owner.clone(),
        total_students: 0,
        total_institutions: 0,
        total_progress_records: 0,
        analytics_enabled: msg.analytics_enabled,
    };

    let config = ProgressConfig {
        update_frequency: msg.update_frequency,
        analytics_depth: msg.analytics_depth,
        retention_period_days: 365, // Default 1 year retention
        dashboard_refresh_hours: 24, // Daily refresh
        benchmark_computation_enabled: true,
    };

    STATE.save(deps.storage, &state)?;
    PROGRESS_CONFIG.save(deps.storage, &config)?;

    Ok(Response::new()
        .add_attribute("method", "instantiate")
        .add_attribute("owner", owner.to_string())
        .add_attribute("analytics_enabled", msg.analytics_enabled.to_string()))
}

#[entry_point]
pub fn execute(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    msg: ExecuteMsg,
) -> Result<Response, ContractError> {
    match msg {
        ExecuteMsg::UpdateConfig { 
            new_owner, 
            analytics_enabled, 
            update_frequency, 
            analytics_depth, 
            retention_period_days 
        } => execute::update_config(deps, info, new_owner, analytics_enabled, update_frequency, analytics_depth, retention_period_days),
        
        ExecuteMsg::UpdateStudentProgress { 
            student_progress, 
            force_analytics_refresh 
        } => execute::update_student_progress(deps, env, student_progress, force_analytics_refresh),
        
        ExecuteMsg::BatchUpdateStudentProgress { 
            student_progress_list, 
            analytics_refresh_mode 
        } => execute::batch_update_student_progress(deps, env, student_progress_list, analytics_refresh_mode),
        
        ExecuteMsg::RecordSubjectCompletion { 
            student_id, 
            subject_id, 
            final_grade, 
            completion_date, 
            study_hours, 
            difficulty_rating, 
            satisfaction_rating, 
            nft_token_id 
        } => execute::record_subject_completion(deps, env, student_id, subject_id, final_grade, completion_date, study_hours, difficulty_rating, satisfaction_rating, nft_token_id),
        
        ExecuteMsg::RecordSubjectEnrollment { 
            student_id, 
            subject_id, 
            enrollment_date, 
            expected_completion 
        } => execute::record_subject_enrollment(deps, student_id, subject_id, enrollment_date, expected_completion),
        
        ExecuteMsg::UpdateCurrentSubjectProgress { 
            student_id, 
            subject_id, 
            current_grade, 
            attendance_rate, 
            assignments_completed, 
            study_hours_logged 
        } => execute::update_current_subject_progress(deps, env, student_id, subject_id, current_grade, attendance_rate, assignments_completed, study_hours_logged),
        
        ExecuteMsg::RecordMilestoneAchievement { 
            student_id, 
            milestone_id, 
            achievement_date, 
            notes 
        } => execute::record_milestone_achievement(deps, env, student_id, milestone_id, achievement_date, notes),
        
        ExecuteMsg::AddRiskFactor { 
            student_id, 
            risk_type, 
            severity, 
            description, 
            intervention_recommendations 
        } => execute::add_risk_factor(deps, env, student_id, risk_type, severity, description, intervention_recommendations),
        
        ExecuteMsg::RemoveRiskFactor { 
            student_id, 
            risk_type 
        } => execute::remove_risk_factor(deps, student_id, risk_type),
        
        ExecuteMsg::GenerateInstitutionAnalytics { 
            institution_id, 
            period, 
            force_refresh 
        } => execute::generate_institution_analytics(deps, env, institution_id, period, force_refresh),
        
        ExecuteMsg::GenerateStudentDashboard { 
            student_id, 
            include_peer_comparison, 
            force_refresh 
        } => execute::generate_student_dashboard(deps, env, student_id, include_peer_comparison, force_refresh),
        
        ExecuteMsg::BatchGenerateDashboards { 
            institution_id, 
            student_ids, 
            include_analytics 
        } => execute::batch_generate_dashboards(deps, env, institution_id, student_ids, include_analytics),
        
        ExecuteMsg::RunGlobalAnalytics { 
            period, 
            institutions 
        } => execute::run_global_analytics(deps, env, period, institutions),
        
        ExecuteMsg::UpdateAcademicStatus { 
            student_id, 
            new_status, 
            effective_date, 
            reason 
        } => execute::update_academic_status(deps, env, student_id, new_status, effective_date, reason),
        
        ExecuteMsg::ArchiveOldData { 
            cutoff_date, 
            dry_run 
        } => execute::archive_old_data(deps, cutoff_date, dry_run),
        
        ExecuteMsg::ResetStudentProgress { 
            student_id, 
            reason, 
            backup 
        } => execute::reset_student_progress(deps, info, student_id, reason, backup),
    }
}

#[entry_point]
pub fn query(deps: Deps, env: Env, msg: QueryMsg) -> StdResult<Binary> {
    match msg {
        QueryMsg::GetState {} => to_json_binary(&query::get_state(deps)?),
        QueryMsg::GetConfig {} => to_json_binary(&query::get_config(deps)?),
        QueryMsg::GetStudentProgress { student_id } => to_json_binary(&query::get_student_progress(deps, student_id)?),
        QueryMsg::GetStudentDashboard { student_id, include_peer_comparison } => to_json_binary(&query::get_student_dashboard(deps, student_id, include_peer_comparison)?),
        QueryMsg::GetInstitutionAnalytics { institution_id, period } => to_json_binary(&query::get_institution_analytics(deps, institution_id, period)?),
        QueryMsg::GetStudentsByInstitution { institution_id, status_filter, limit, start_after } => to_json_binary(&query::get_students_by_institution(deps, institution_id, status_filter, limit, start_after)?),
        QueryMsg::GetStudentsByCourse { course_id, status_filter, limit, start_after } => to_json_binary(&query::get_students_by_course(deps, course_id, status_filter, limit, start_after)?),
        QueryMsg::GetAtRiskStudents { institution_id, risk_level_filter, limit } => to_json_binary(&query::get_at_risk_students(deps, institution_id, risk_level_filter, limit)?),
        QueryMsg::GetTopPerformers { institution_id, course_id, metric, limit } => to_json_binary(&query::get_top_performers(deps, institution_id, course_id, metric, limit)?),
        QueryMsg::GetGraduationForecast { student_id, scenario } => to_json_binary(&query::get_graduation_forecast(deps, student_id, scenario)?),
        QueryMsg::GetComparativeAnalytics { institution_id, comparison_group, metrics } => to_json_binary(&query::get_comparative_analytics(deps, institution_id, comparison_group, metrics)?),
        QueryMsg::GetProgressTrends { entity_id, entity_type, period, metrics } => to_json_binary(&query::get_progress_trends(deps, entity_id, entity_type, period, metrics)?),
        QueryMsg::GetSubjectPerformance { subject_id, institution_id, period } => to_json_binary(&query::get_subject_performance(deps, subject_id, institution_id, period)?),
        QueryMsg::GetCohortAnalysis { institution_id, course_id, enrollment_year, analysis_type } => to_json_binary(&query::get_cohort_analysis(deps, institution_id, course_id, enrollment_year, analysis_type)?),
        QueryMsg::GetPredictiveInsights { student_id, prediction_horizon_semesters } => to_json_binary(&query::get_predictive_insights(deps, student_id, prediction_horizon_semesters)?),
        QueryMsg::SearchStudents { query, limit, start_after } => to_json_binary(&query::search_students(deps, query, limit, start_after)?),
        QueryMsg::GetAnalyticsSummary { scope, quick_stats_only } => to_json_binary(&query::get_analytics_summary(deps, env, scope, quick_stats_only)?),
    }
}

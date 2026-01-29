import React, { useState, useEffect, useMemo } from 'react';

const API_URL = 'http://localhost:8080/api';

// 日历组件
function Calendar({ records, onDateClick, selectedDate }) {
  const [currentMonth, setCurrentMonth] = useState(new Date());

  const year = currentMonth.getFullYear();
  const month = currentMonth.getMonth();

  const daysInMonth = new Date(year, month + 1, 0).getDate();
  const firstDayOfMonth = new Date(year, month, 1).getDay();

  const recordsByDate = useMemo(() => {
    const map = {};
    records.forEach(r => {
      if (!map[r.date]) map[r.date] = [];
      map[r.date].push(r);
    });
    return map;
  }, [records]);

  const days = [];
  for (let i = 0; i < firstDayOfMonth; i++) {
    days.push(<div key={`empty-${i}`} className="calendar-day empty"></div>);
  }

  for (let day = 1; day <= daysInMonth; day++) {
    const dateStr = `${year}-${String(month + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`;
    const hasRecords = recordsByDate[dateStr];
    const isSelected = selectedDate === dateStr;
    const isToday = dateStr === new Date().toISOString().split('T')[0];

    days.push(
      <div
        key={day}
        className={`calendar-day ${hasRecords ? 'has-record' : ''} ${isSelected ? 'selected' : ''} ${isToday ? 'today' : ''}`}
        onClick={() => onDateClick(dateStr)}
      >
        <span className="day-number">{day}</span>
        {hasRecords && <span className="record-count">{hasRecords.length}</span>}
      </div>
    );
  }

  const monthNames = ['一月', '二月', '三月', '四月', '五月', '六月', '七月', '八月', '九月', '十月', '十一月', '十二月'];
  const weekDays = ['日', '一', '二', '三', '四', '五', '六'];

  return (
    <div className="calendar">
      <div className="calendar-header">
        <button onClick={() => setCurrentMonth(new Date(year, month - 1))}>◀</button>
        <h3>{year}年 {monthNames[month]}</h3>
        <button onClick={() => setCurrentMonth(new Date(year, month + 1))}>▶</button>
      </div>
      <div className="calendar-weekdays">
        {weekDays.map(d => <div key={d} className="weekday">{d}</div>)}
      </div>
      <div className="calendar-grid">{days}</div>
    </div>
  );
}

// 频率热力图组件（类似GitHub贡献图）
function FrequencyChart({ records }) {
  const weeks = 52;
  const today = new Date();

  const recordsByDate = useMemo(() => {
    const map = {};
    records.forEach(r => {
      map[r.date] = (map[r.date] || 0) + 1;
    });
    return map;
  }, [records]);

  const maxCount = Math.max(1, ...Object.values(recordsByDate));

  const getIntensity = (count) => {
    if (count === 0) return 0;
    if (count <= maxCount * 0.25) return 1;
    if (count <= maxCount * 0.5) return 2;
    if (count <= maxCount * 0.75) return 3;
    return 4;
  };

  const cells = [];
  const startDate = new Date(today);
  startDate.setDate(startDate.getDate() - (weeks * 7) + 1);
  startDate.setDate(startDate.getDate() - startDate.getDay());

  for (let week = 0; week < weeks; week++) {
    const weekCells = [];
    for (let day = 0; day < 7; day++) {
      const cellDate = new Date(startDate);
      cellDate.setDate(cellDate.getDate() + week * 7 + day);
      const dateStr = cellDate.toISOString().split('T')[0];
      const count = recordsByDate[dateStr] || 0;
      const intensity = getIntensity(count);

      weekCells.push(
        <div
          key={`${week}-${day}`}
          className={`freq-cell intensity-${intensity}`}
          title={`${dateStr}: ${count}次`}
        />
      );
    }
    cells.push(<div key={week} className="freq-week">{weekCells}</div>);
  }

  const monthLabels = [];
  let lastMonth = -1;
  for (let week = 0; week < weeks; week++) {
    const cellDate = new Date(startDate);
    cellDate.setDate(cellDate.getDate() + week * 7);
    const month = cellDate.getMonth();
    if (month !== lastMonth) {
      const monthNames = ['1月', '2月', '3月', '4月', '5月', '6月', '7月', '8月', '9月', '10月', '11月', '12月'];
      monthLabels.push(
        <span key={week} style={{ left: `${(week / weeks) * 100}%` }} className="month-label">
          {monthNames[month]}
        </span>
      );
      lastMonth = month;
    }
  }

  return (
    <div className="frequency-chart">
      <h2>活动频率</h2>
      <div className="freq-container">
        <div className="freq-months">{monthLabels}</div>
        <div className="freq-grid">{cells}</div>
        <div className="freq-legend">
          <span>少</span>
          <div className="freq-cell intensity-0" />
          <div className="freq-cell intensity-1" />
          <div className="freq-cell intensity-2" />
          <div className="freq-cell intensity-3" />
          <div className="freq-cell intensity-4" />
          <span>多</span>
        </div>
      </div>
    </div>
  );
}

function App() {
  const [records, setRecords] = useState([]);
  const [stats, setStats] = useState({ totalRecords: 0, totalDuration: 0 });
  const [form, setForm] = useState({
    date: new Date().toISOString().split('T')[0],
    content: '',
    duration: '',
    notes: ''
  });
  const [editingId, setEditingId] = useState(null);
  const [selectedDate, setSelectedDate] = useState(null);

  const fetchRecords = async () => {
    try {
      const res = await fetch(`${API_URL}/records`);
      const data = await res.json();
      setRecords(data || []);
    } catch (err) {
      console.error('Failed to fetch records:', err);
    }
  };

  const fetchStats = async () => {
    try {
      const res = await fetch(`${API_URL}/stats`);
      const data = await res.json();
      setStats(data);
    } catch (err) {
      console.error('Failed to fetch stats:', err);
    }
  };

  useEffect(() => {
    fetchRecords();
    fetchStats();
  }, []);

  const handleSubmit = async (e) => {
    e.preventDefault();
    const payload = {
      ...form,
      duration: parseInt(form.duration) || 0
    };

    try {
      if (editingId) {
        await fetch(`${API_URL}/records/${editingId}`, {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(payload)
        });
        setEditingId(null);
      } else {
        await fetch(`${API_URL}/records`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(payload)
        });
      }
      setForm({
        date: new Date().toISOString().split('T')[0],
        content: '',
        duration: '',
        notes: ''
      });
      fetchRecords();
      fetchStats();
    } catch (err) {
      console.error('Failed to save record:', err);
    }
  };

  const handleDelete = async (id) => {
    if (!window.confirm('确定要删除这条记录吗？')) return;
    try {
      await fetch(`${API_URL}/records/${id}`, { method: 'DELETE' });
      fetchRecords();
      fetchStats();
    } catch (err) {
      console.error('Failed to delete record:', err);
    }
  };

  const handleEdit = (record) => {
    setForm({
      date: record.date,
      content: record.content,
      duration: record.duration.toString(),
      notes: record.notes
    });
    setEditingId(record.id);
  };

  const cancelEdit = () => {
    setEditingId(null);
    setForm({
      date: new Date().toISOString().split('T')[0],
      content: '',
      duration: '',
      notes: ''
    });
  };

  const handleDateClick = (dateStr) => {
    setSelectedDate(selectedDate === dateStr ? null : dateStr);
    setForm(prev => ({ ...prev, date: dateStr }));
  };

  const formatDuration = (minutes) => {
    if (minutes < 60) return `${minutes}分钟`;
    const hours = Math.floor(minutes / 60);
    const mins = minutes % 60;
    return mins > 0 ? `${hours}小时${mins}分钟` : `${hours}小时`;
  };

  const filteredRecords = selectedDate
    ? records.filter(r => r.date === selectedDate)
    : records;

  return (
    <div className="container">
      <h1>个人习惯追踪</h1>

      <div className="stats">
        <div className="stat-card">
          <h3>总记录数</h3>
          <div className="value">{stats.totalRecords}</div>
        </div>
        <div className="stat-card">
          <h3>总时长</h3>
          <div className="value">{formatDuration(stats.totalDuration)}</div>
        </div>
      </div>

      <FrequencyChart records={records} />

      <div className="main-content">
        <div className="left-panel">
          <Calendar
            records={records}
            onDateClick={handleDateClick}
            selectedDate={selectedDate}
          />
        </div>

        <div className="right-panel">
          <div className="form-section">
            <h2>{editingId ? '编辑记录' : '添加新记录'}</h2>
            <form onSubmit={handleSubmit}>
              <div className="form-row">
                <div className="form-group">
                  <label>日期</label>
                  <input
                    type="date"
                    value={form.date}
                    onChange={(e) => setForm({ ...form, date: e.target.value })}
                    required
                  />
                </div>
                <div className="form-group">
                  <label>时长（分钟）</label>
                  <input
                    type="number"
                    value={form.duration}
                    onChange={(e) => setForm({ ...form, duration: e.target.value })}
                    placeholder="例如: 30"
                    min="1"
                    required
                  />
                </div>
              </div>
              <div className="form-group">
                <label>观看内容</label>
                <input
                  type="text"
                  value={form.content}
                  onChange={(e) => setForm({ ...form, content: e.target.value })}
                  placeholder="输入观看的内容名称"
                  required
                />
              </div>
              <div className="form-group">
                <label>备注</label>
                <textarea
                  value={form.notes}
                  onChange={(e) => setForm({ ...form, notes: e.target.value })}
                  placeholder="可选备注..."
                />
              </div>
              <button type="submit" className="btn btn-primary">
                {editingId ? '更新' : '保存'}
              </button>
              {editingId && (
                <button type="button" className="btn btn-cancel" onClick={cancelEdit}>
                  取消
                </button>
              )}
            </form>
          </div>
        </div>
      </div>

      <div className="records-section">
        <h2>
          {selectedDate ? `${selectedDate} 的记录` : '历史记录'}
          {selectedDate && (
            <button className="btn-clear" onClick={() => setSelectedDate(null)}>
              显示全部
            </button>
          )}
        </h2>
        {filteredRecords.length === 0 ? (
          <div className="empty-state">暂无记录</div>
        ) : (
          <div className="record-list">
            {filteredRecords.sort((a, b) => new Date(b.date) - new Date(a.date)).map((record) => (
              <div key={record.id} className="record-card">
                <div className="record-info">
                  <h3>{record.content}</h3>
                  <p>日期: {record.date}</p>
                  <p>时长: {formatDuration(record.duration)}</p>
                  {record.notes && <p className="notes">备注: {record.notes}</p>}
                </div>
                <div className="record-actions">
                  <button className="btn btn-primary" onClick={() => handleEdit(record)}>
                    编辑
                  </button>
                  <button className="btn btn-danger" onClick={() => handleDelete(record.id)}>
                    删除
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

export default App;

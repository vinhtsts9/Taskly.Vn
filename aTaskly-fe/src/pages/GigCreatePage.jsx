import React, { useContext, useEffect, useMemo, useRef, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { AuthContext } from '../context/AuthContext';
import { apiGetAuth, apiPostAuth, apiPutAuth, apiUploadAuth } from '../utils/api';
import './DashboardPage.css';

const Stepper = ({ step, total = 5 }) => (
  <div style={{ display: 'flex', gap: 8, marginBottom: 16 }}>
    {Array.from({ length: total }).map((_, i) => (
      <div key={i} style={{
        flex: 1, height: 6, borderRadius: 4,
        background: i < step ? '#1dbf73' : '#e5e7eb'
      }} />
    ))}
  </div>
);

const FilePicker = ({ label, multiple = true, accept, maxSizeMB = 5, onUpload }) => {
  const inputRef = useRef(null);
  const onChange = async (e) => {
    const files = Array.from(e.target.files || []);
    const maxBytes = maxSizeMB * 1024 * 1024;
    const valid = files.filter(f => f.size <= maxBytes);
    if (valid.length !== files.length) {
      alert(`Một số file vượt quá ${maxSizeMB}MB và đã bị loại bỏ.`);
    }
    if (valid.length === 0) return;
    const form = new FormData();
    valid.forEach((f) => form.append('files', f));
    await onUpload(form);
    if (inputRef.current) inputRef.current.value = '';
  };
  return (
    <div>
      <div style={{ fontSize: 12, color: '#6b7280', marginBottom: 6 }}>{label}</div>
      <input ref={inputRef} type="file" multiple={multiple} accept={accept} onChange={onChange} />
      <div style={{ fontSize: 12, color: '#94a3b8' }}>Giới hạn {maxSizeMB}MB mỗi file</div>
    </div>
  );
};

const GigCreatePage = () => {
  const navigate = useNavigate();
  const { isAuthenticated, currentUser } = useContext(AuthContext);

  const [step, setStep] = useState(1);
  const [loading, setLoading] = useState(false);
  const [msg, setMsg] = useState('');

  // Step 1
  const [title, setTitle] = useState('');
  const [categories, setCategories] = useState([]);
  const [categoryId, setCategoryId] = useState('');
  const [subcategoryId, setSubcategoryId] = useState('');

  // Step 2
  const [description, setDescription] = useState('');
  const [images, setImages] = useState([]); // urls

  // Step 3
  const [pricingMode, setPricingMode] = useState('single'); // 'single' | 'triple'
  const [packages, setPackages] = useState([
    { tier: 'basic',  price: '', delivery_days: '', options: { revisions: 1, files: 1, extras: [] } }
  ]);

  // Step 4
  const [questions, setQuestions] = useState([]);

  useEffect(() => {
    if (!isAuthenticated) {
      navigate('/login');
    }
  }, [isAuthenticated, navigate]);

  useEffect(() => {
    const fetchCategories = async () => {
      try {
        const res = await apiGetAuth('/gigs/categories');
        // res là mảng flat
        const rawCategories = Array.isArray(res) ? res : (res?.data || []);
  
        // Tạo map cha -> object { id, name, children }
        const categoriesMap = new Map();
  
        rawCategories.forEach(({ parent_id, parent_name, children_id, children_name }) => {
          if (!categoriesMap.has(parent_id)) {
            categoriesMap.set(parent_id, { id: parent_id, name: parent_name, children: [] });
          }
          if (children_id !== null) {
            categoriesMap.get(parent_id).children.push({ id: children_id, name: children_name });
          }
        });
  
        // Chuyển map thành mảng
        const categoriesTree = Array.from(categoriesMap.values());
  
        setCategories(categoriesTree);
      } catch (e) {
        setCategories([]);
      }
    };
    fetchCategories();
  }, []);
  

  const selectedCategory = categories.find(c => String(c.id) === String(categoryId));
  const subcategories = selectedCategory?.subcategories || selectedCategory?.children || [];

  const toTriple = () => {
    setPricingMode('triple');
    setPackages([
      { tier: 'basic', price: '', delivery_days: '', options: { revisions: 1, files: 1, extras: [] } },
      { tier: 'standard', price: '', delivery_days: '', options: { revisions: 1, files: 1, extras: [] } },
      { tier: 'premium', price: '', delivery_days: '', options: { revisions: 1, files: 1, extras: [] } },
    ]);
  };
  const toSingle = () => {
    setPricingMode('single');
    setPackages([
      { tier: 'basic', price: '', delivery_days: '', options: { revisions: 1, files: 1, extras: [] } },
    ]);
  };

  const setPackageField = (idx, field, value) => {
    setPackages((prev) => prev.map((p, i) => i === idx ? { ...p, [field]: value } : p));
  };

  const addQuestion = () => setQuestions((qs) => [...qs, { question: '', type: 'text', required: true, options: [] }]);
  const removeQuestion = (i) => setQuestions((qs) => qs.filter((_, idx) => idx !== i));
  const setQuestionField = (i, field, value) => setQuestions((qs) => qs.map((q, idx) => idx === i ? { ...q, [field]: value } : q));

  const uploadImages = async (form) => {
    const res = await apiUploadAuth('/gigs/upload-media', form);
    const urls = res?.urls || res || [];
    setImages((prev) => [...prev, ...urls]);
  };

  const validateStep1 = () => title.trim().length >= 10 && categoryId;
  const validateStep2 = () => description.trim().length >= 30; // có media hay không tùy backend
  const validateStep3 = () => packages.every(p => Number(p.price) > 0 && Number(p.delivery_days) >= 1);

  const handlePublish = async (status) => {
    // Lấy ID danh mục và danh mục con, đảm bảo là số nguyên và tạo thành mảng
    const gigCategoryIds = [];
    if (categoryId) {
        gigCategoryIds.push(parseInt(categoryId, 10)); // Chuyển từ string UUID sang int
    }
    if (subcategoryId) {
        gigCategoryIds.push(parseInt(subcategoryId, 10)); // Chuyển từ string UUID sang int
    }

    // Kiểm tra currentUser trước khi sử dụng
    if (!currentUser || !currentUser.id) {
      setMsg('Lỗi: Không thể lấy thông tin người dùng. Vui lòng đăng nhập lại.');
      return;
    }

    const payload = {
      user_id: currentUser.id, // Thay user.id bằng currentUser.id
      title,
      category_id: gigCategoryIds.length > 0 ? gigCategoryIds : [], // Gửi mảng INT[]
      description,
      image_url: images, // Chỉ gửi mảng ảnh, bỏ qua video
      pricing_mode: pricingMode,
      packages: packages.map(pkg => ({
        tier: pkg.tier,
        price: Number(pkg.price), // Đảm bảo là số
        delivery_days: Number(pkg.delivery_days), // Đảm bảo là số
        options: { // Bao gồm options
          revisions: Number(pkg.options.revisions),
          files: Number(pkg.options.files),
        }
      })),
      requirements: { 
        questions: questions.map(q => ({
          question: q.question,
          required: q.required,
          // Bỏ qua type và options vì backend gig_requirements không có cột này
        }))
      },
      status,
    };
    setLoading(true);
    setMsg('');
    try {
      console.log("This is payload:  ",payload)
      await apiPostAuth('/gigs', payload);
      navigate('/my-gigs'); // Chuyển hướng đến trang /my-gigs bất kể trạng thái
      setMsg(status === 'draft' ? 'Đã lưu nháp thành công!' : 'Gig đã được xuất bản thành công!');
    } catch (e) {
      setMsg(e?.message || 'Không thể lưu gig');
    } finally {
      setLoading(false);
    }
  };

  const canNext = useMemo(() => {
    if (step === 1) return validateStep1();
    if (step === 2) return validateStep2();
    if (step === 3) return validateStep3();
    return true;
  }, [step, title, categoryId, description, packages]);

  return (
    <div className="dashboard-page">
      <div className="dashboard-header">
        <h1>Đăng gig mới</h1>
        <h2>5 bước đơn giản để xuất bản dịch vụ của bạn</h2>
      </div>
      <Stepper step={step} total={5} />
      {msg && <div style={{ marginBottom: 8 }}>{msg}</div>}
      {loading && <div style={{ marginBottom: 8 }}>Đang xử lý...</div>}

      {step === 1 && (
        <div className="orders-table-container">
          <div className="admin-row" style={{ gap: 16 }}>
            <div className="admin-field" style={{ flex: 1 }}>
              <div className="admin-field-label">Tiêu đề</div>
              <input value={title} onChange={(e) => setTitle(e.target.value)} placeholder="VD: Thiết kế logo thương hiệu trong 3 ngày" />
              <div className="admin-help">Tối thiểu 10 ký tự</div>
            </div>
          </div>
          <div className="admin-row" style={{ gap: 16, marginTop: 12 }}>
            <div className="admin-field" style={{ flex: 1 }}>
              <div className="admin-field-label">Danh mục</div>
              <select value={categoryId} onChange={(e) => { setCategoryId(e.target.value); setSubcategoryId(''); }}>
                <option value="">-- Chọn danh mục --</option>
                {categories.map(c => <option key={c.id} value={c.id}>{c.name || c.title}</option>)}
              </select>
            </div>
            <div className="admin-field" style={{ flex: 1 }}>
              <div className="admin-field-label">Danh mục con</div>
              <select value={subcategoryId} onChange={(e) => setSubcategoryId(e.target.value)} disabled={!subcategories?.length}>
                <option value="">-- Không chọn --</option>
                {subcategories?.map(sc => <option key={sc.id} value={sc.id}>{sc.name || sc.title}</option>)}
              </select>
              <div className="admin-help">Dữ liệu tải từ backend</div>
            </div>
          </div>
        </div>
      )}

      {step === 2 && (
        <div className="orders-table-container">
          <div className="admin-field">
            <div className="admin-field-label">Mô tả dịch vụ</div>
            <textarea rows={8} value={description} onChange={(e) => setDescription(e.target.value)} placeholder="Hãy mô tả rõ phạm vi công việc, quy trình, những gì bạn cần từ khách hàng..." />
            <div className="admin-help">Tối thiểu 30 ký tự</div>
          </div>
          <div className="admin-row" style={{ gap: 16, marginTop: 12 }}>
            <FilePicker label="Ảnh giới thiệu (tối đa 5MB/ảnh)" multiple accept="image/*" maxSizeMB={5} onUpload={uploadImages} />
          </div>
          <div className="admin-row" style={{ gap: 10, marginTop: 12, flexWrap: 'wrap' }}>
            {images.map((url) => (
              <img key={url} src={url} alt="preview" style={{ width: 100, height: 70, objectFit: 'cover', borderRadius: 6, border: '1px solid #eee' }} />
            ))}
            {/* {video && (
              <video src={video} controls style={{ width: 200, borderRadius: 6, border: '1px solid #eee' }} />
            )} */}
          </div>
        </div>
      )}

      {step === 3 && (
        <div className="orders-table-container">
          <div className="admin-row" style={{ justifyContent: 'space-between' }}>
            <div className="admin-help">Chọn chế độ gói</div>
            <div style={{ display: 'flex', gap: 8 }}>
              <button className="admin-button secondary" onClick={toSingle} disabled={pricingMode === 'single'}>1 gói</button>
              <button className="admin-button secondary" onClick={toTriple} disabled={pricingMode === 'triple'}>3 gói</button>
            </div>
          </div>
          <div className="admin-row" style={{ gap: 16, flexWrap: 'wrap' }}>
            {packages.map((pkg, idx) => (
              <div key={pkg.tier} className="admin-card" style={{ flex: '1 1 300px' }}>
                <div className="admin-card-title" style={{ textTransform: 'capitalize' }}>{pkg.tier}</div>
                <div className="admin-row" style={{ gap: 12 }}>
                  <div className="admin-field" style={{ minWidth: 140 }}>
                    <div className="admin-field-label">Giá (VND)</div>
                    <input type="number" min={0} value={pkg.price} onChange={(e) => setPackageField(idx, 'price', e.target.value)} />
                  </div>
                  <div className="admin-field" style={{ minWidth: 140 }}>
                    <div className="admin-field-label">Thời gian (ngày)</div>
                    <input type="number" min={1} value={pkg.delivery_days} onChange={(e) => setPackageField(idx, 'delivery_days', e.target.value)} />
                  </div>
                </div>
                <div className="admin-row" style={{ gap: 12 }}>
                  <div className="admin-field" style={{ minWidth: 140 }}>
                    <div className="admin-field-label">Số lần sửa</div>
                    <input type="number" min={0} value={pkg.options.revisions} onChange={(e) => setPackageField(idx, 'options', { ...pkg.options, revisions: Number(e.target.value) })} />
                  </div>
                  <div className="admin-field" style={{ minWidth: 140 }}>
                    <div className="admin-field-label">Số file bàn giao</div>
                    <input type="number" min={1} value={pkg.options.files} onChange={(e) => setPackageField(idx, 'options', { ...pkg.options, files: Number(e.target.value) })} />
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {step === 4 && (
        <div className="orders-table-container">
          <div className="admin-row" style={{ justifyContent: 'space-between' }}>
            <div className="admin-card-title">Câu hỏi cho người mua</div>
            <button className="admin-button" onClick={addQuestion}>Thêm câu hỏi</button>
          </div>
          {questions.length === 0 && <div className="admin-help">Chưa có câu hỏi nào</div>}
          <div className="admin-row" style={{ gap: 16, flexDirection: 'column', alignItems: 'stretch' }}>
            {questions.map((q, i) => (
              <div key={i} className="admin-card">
                <div className="admin-row" style={{ justifyContent: 'space-between' }}>
                  <div className="admin-card-title">Câu hỏi #{i + 1}</div>
                  <button className="admin-button secondary" onClick={() => removeQuestion(i)}>Xóa</button>
                </div>
                <div className="admin-field">
                  <div className="admin-field-label">Nội dung câu hỏi</div>
                  <input value={q.question} onChange={(e) => setQuestionField(i, 'question', e.target.value)} />
                </div>
                <div className="admin-row" style={{ gap: 12 }}>
                  <div className="admin-field">
                    <div className="admin-field-label">Bắt buộc</div>
                    <select value={q.required ? '1' : '0'} onChange={(e) => setQuestionField(i, 'required', e.target.value === '1')}>
                      <option value="1">Có</option>
                      <option value="0">Không</option>
                    </select>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {step === 5 && (
        <div className="orders-table-container">
          <div className="admin-card-title">Xuất bản</div>
          <div className="admin-help">Bạn có thể lưu nháp để hoàn thiện sau hoặc xuất bản ngay khi nội dung đã đầy đủ.</div>
          <div style={{ display: 'flex', gap: 8, marginTop: 12 }}>
            <button className="admin-button secondary" onClick={() => handlePublish('draft')}>Lưu nháp</button>
            <button className="admin-button" onClick={() => handlePublish('active')}>Xuất bản</button>
          </div>
        </div>
      )}

      <div style={{ display: 'flex', justifyContent: 'space-between', marginTop: 16 }}>
        <button className="admin-button secondary" disabled={step === 1} onClick={() => setStep((s) => Math.max(1, s - 1))}>Quay lại</button>
        <button className="admin-button" disabled={!canNext || step === 5} onClick={() => setStep((s) => Math.min(5, s + 1))}>Tiếp tục</button>
      </div>
    </div>
  );
};

export default GigCreatePage;



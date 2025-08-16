import React, { useEffect, useMemo, useState } from 'react';
import { apiGetAuth, apiPostAuth, apiDeleteAuthWithBody } from '../utils/api';
import './AdminDashboardPage.css';

const Section = ({ title, subtitle, children }) => (
  <section className="admin-card">
    <div className="admin-card-title">{title}</div>
    {subtitle && <div className="admin-help" style={{ marginBottom: 8 }}>{subtitle}</div>}
    {children}
  </section>
);

const Field = ({ label, children }) => (
  <label className="admin-field">
    <div className="admin-field-label">{label}</div>
    {children}
  </label>
);

const Row = ({ children }) => (
  <div className="admin-row" style={{ marginBottom: 8 }}>
    {children}
  </div>
);

const Button = ({ children, variant = 'primary', ...props }) => (
  <button {...props} className={`admin-button ${variant}`}>{children}</button>
);

const Code = ({ children }) => (
  <pre className="admin-code">
    <code>{children}</code>
  </pre>
);

const AdminDashboardPage = () => {
  // Form states
  const [roleName, setRoleName] = useState('');
  const [permission, setPermission] = useState({ name: '', resource: '', action: '' });
  const [userRole, setUserRole] = useState({ user_id: '', role_id: '' });
  const [rolePerm, setRolePerm] = useState({ role_id: '', permission_id: '' });
  const [query, setQuery] = useState({ user_id: '', role_id: '' });

  // Data states
  const [userPermissions, setUserPermissions] = useState([]);
  const [userRoles, setUserRoles] = useState([]);
  const [rolePermissions, setRolePermissions] = useState([]);
  const [busy, setBusy] = useState(false);
  const [msg, setMsg] = useState('');

  // Users tab states
  const [users, setUsers] = useState([]);
  const [usersTotal, setUsersTotal] = useState(0);
  const [userSearch, setUserSearch] = useState('');
  const [page, setPage] = useState(1);
  const [size, setSize] = useState(10);
  const [usersLoading, setUsersLoading] = useState(false);
  const [selectedUserId, setSelectedUserId] = useState('');
  const [userDetail, setUserDetail] = useState(null);
  const [userDetailLoading, setUserDetailLoading] = useState(false);
  const [userDetailError, setUserDetailError] = useState('');

  // Helpers for formatting
  const maskEmail = (email) => {
    if (!email) return '';
    const str = String(email);
    if (str.length <= 7) return str;
    return `${str.slice(0, 2)}***${str.slice(-5)}`;
  };

  const formatDate = (isoOrPgTimestamp) => {
    if (!isoOrPgTimestamp) return '';
    const d = new Date(isoOrPgTimestamp);
    if (Number.isNaN(d.getTime())) return String(isoOrPgTimestamp);
    return d.toLocaleString('vi-VN', {
      year: 'numeric', month: '2-digit', day: '2-digit',
      hour: '2-digit', minute: '2-digit', second: '2-digit'
    });
  };

  const getStatePill = (states) => {
    const v = typeof states === 'string' ? states.trim() : states;
    if (v === 1 || v === '1' || v === true || v === 'active') {
      return { label: 'Hoạt động', cls: 'admin-pill success' };
    }
    if (v === 0 || v === '0' || v === false || v === 'disabled') {
      return { label: 'Bị khóa', cls: 'admin-pill warn' };
    }
    return { label: 'Không rõ', cls: 'admin-pill muted' };
  };

  const notEmpty = (v) => v && String(v).trim().length > 0;
  const canCreateRole = useMemo(() => notEmpty(roleName), [roleName]);
  const canCreatePermission = useMemo(() => notEmpty(permission.name) && notEmpty(permission.resource) && notEmpty(permission.action), [permission]);
  const canMapUserRole = useMemo(() => notEmpty(userRole.user_id) && notEmpty(userRole.role_id), [userRole]);
  const canMapRolePerm = useMemo(() => notEmpty(rolePerm.role_id) && notEmpty(rolePerm.permission_id), [rolePerm]);

  const call = async (fn, successMsg) => {
    try {
      setBusy(true);
      setMsg('');
      await fn();
      setMsg(successMsg);
    } catch (e) {
      setMsg(e?.message || 'Có lỗi xảy ra');
    } finally {
      setBusy(false);
    }
  };

  // Actions mapping to backend routes
  const createRole = () => call(() => apiPostAuth('/admin/rbac/role', { name: roleName }), 'Tạo role thành công');
  const createPermission = () => call(() => apiPostAuth('/admin/rbac/permission', permission), 'Tạo permission thành công');
  const addRoleToUser = () => call(() => apiPostAuth('/admin/rbac/user-role', userRole), 'Gán role cho user thành công');
  const removeRoleFromUser = () => call(() => apiDeleteAuthWithBody('/admin/rbac/user-role', userRole), 'Hủy gán role khỏi user thành công');
  const addPermissionToRole = () => call(() => apiPostAuth('/admin/rbac/role-permission', rolePerm), 'Gán permission cho role thành công');
  const removePermissionFromRole = () => call(() => apiDeleteAuthWithBody('/admin/rbac/role-permission', rolePerm), 'Hủy gán permission khỏi role thành công');
  const fetchPermissionsByUser = () => call(async () => setUserPermissions(await apiGetAuth(`/admin/rbac/roles/${encodeURIComponent(query.user_id)}/permissions`)), 'Lấy permissions theo user thành công');
  const fetchRolesByUser = () => call(async () => setUserRoles(await apiGetAuth(`/admin/rbac/roles/${encodeURIComponent(query.user_id)}`)), 'Lấy roles theo user thành công');
  const fetchPermissionsByRole = () => call(async () => setRolePermissions(await apiGetAuth(`/admin/rbac/permissions/${encodeURIComponent(query.role_id)}`)), 'Lấy permissions theo role thành công');

  // Users tab actions
  const fetchUsers = async (opts = {}) => {
    const q = opts.query ?? userSearch;
    const p = opts.page ?? page;
    const s = opts.size ?? size;
    try {
      setUsersLoading(true);
      setMsg('');
      const res = await apiGetAuth(`/admin/users?query=${encodeURIComponent(q)}&page=${encodeURIComponent(p)}&size=${encodeURIComponent(s)}`);
      // Backend returns: { total, page, size, data }
      if (Array.isArray(res)) {
        setUsers(res);
        setUsersTotal(res.length);
      } else {
        const items = res.data || [];
        setUsers(items);
        setUsersTotal(res.total ?? (Array.isArray(items) ? items.length : 0));
        if (typeof res.page === 'number') setPage(res.page);
        if (typeof res.size === 'number') setSize(res.size);
      }
    } catch (e) {
      setMsg(e?.message || 'Không thể tải danh sách người dùng');
      setUsers([]);
      setUsersTotal(0);
    } finally {
      setUsersLoading(false);
    }
  };

  useEffect(() => {
    // auto load users list on mount
    fetchUsers().catch(() => {});
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const onSearchUsers = async () => {
    setPage(1);
    await fetchUsers({ query: userSearch, page: 1, size });
  };

  const onChangePage = async (next) => {
    const newPage = Math.max(1, page + next);
    setPage(newPage);
    await fetchUsers({ page: newPage, size, query: userSearch });
  };

  const onChangeSize = async (e) => {
    const newSize = Number(e.target.value) || 10;
    setSize(newSize);
    setPage(1);
    await fetchUsers({ page: 1, size: newSize, query: userSearch });
  };

  const viewUserDetail = async (id) => {
    setSelectedUserId(id);
    setUserDetail(null);
    setUserDetailError('');
    try {
      setUserDetailLoading(true);
      const detail = await apiGetAuth(`/admin/users/${encodeURIComponent(id)}`);
      setUserDetail(detail);
    } catch (e) {
      setUserDetailError(e?.message || 'Không thể tải chi tiết người dùng');
    } finally {
      setUserDetailLoading(false);
    }
  };

  return (
    <div className="admin-page">
      <div className="admin-header">
        <h2>Quản trị RBAC</h2>
        <div className="admin-subtitle">Quản lý vai trò, quyền và người dùng.</div>
      </div>

      {msg && <div className="admin-message">{msg}</div>}
      {busy && <div className="admin-message">Đang xử lý...</div>}

      <Section title="Tạo Role">
        <Row>
          <Field label="Role name">
            <input value={roleName} onChange={(e) => setRoleName(e.target.value)} placeholder="vd: admin" />
          </Field>
          <Button disabled={!canCreateRole || busy} onClick={createRole}>Tạo</Button>
        </Row>
      </Section>

      <Section title="Tạo Permission">
        <Row>
          <Field label="Name">
            <input value={permission.name} onChange={(e) => setPermission({ ...permission, name: e.target.value })} placeholder="vd: manage_users" />
          </Field>
          <Field label="Resource">
            <input value={permission.resource} onChange={(e) => setPermission({ ...permission, resource: e.target.value })} placeholder="vd: users" />
          </Field>
          <Field label="Action">
            <input value={permission.action} onChange={(e) => setPermission({ ...permission, action: e.target.value })} placeholder="vd: read|write|delete" />
          </Field>
          <Button disabled={!canCreatePermission || busy} onClick={createPermission}>Tạo</Button>
        </Row>
      </Section>

      <Section title="Gán/Hủy Role cho User">
        <Row>
          <Field label="User ID (uuid)">
            <input value={userRole.user_id} onChange={(e) => setUserRole({ ...userRole, user_id: e.target.value })} placeholder="uuid user" />
          </Field>
          <Field label="Role ID (uuid)">
            <input value={userRole.role_id} onChange={(e) => setUserRole({ ...userRole, role_id: e.target.value })} placeholder="uuid role" />
          </Field>
          <Button disabled={!canMapUserRole || busy} onClick={addRoleToUser}>Gán</Button>
          <Button variant="secondary" disabled={!canMapUserRole || busy} onClick={removeRoleFromUser}>Hủy gán</Button>
        </Row>
      </Section>

      <Section title="Gán/Hủy Permission cho Role">
        <Row>
          <Field label="Role ID (uuid)">
            <input value={rolePerm.role_id} onChange={(e) => setRolePerm({ ...rolePerm, role_id: e.target.value })} placeholder="uuid role" />
          </Field>
          <Field label="Permission ID (uuid)">
            <input value={rolePerm.permission_id} onChange={(e) => setRolePerm({ ...rolePerm, permission_id: e.target.value })} placeholder="uuid permission" />
          </Field>
          <Button disabled={!canMapRolePerm || busy} onClick={addPermissionToRole}>Gán</Button>
          <Button variant="secondary" disabled={!canMapRolePerm || busy} onClick={removePermissionFromRole}>Hủy gán</Button>
        </Row>
      </Section>

      <Section title="Truy vấn">
        <Row>
          <Field label="User ID (uuid)">
            <input value={query.user_id} onChange={(e) => setQuery({ ...query, user_id: e.target.value })} placeholder="uuid user" />
          </Field>
          <Button disabled={!notEmpty(query.user_id) || busy} onClick={fetchPermissionsByUser}>Lấy Permissions theo User</Button>
          <Button variant="secondary" disabled={!notEmpty(query.user_id) || busy} onClick={fetchRolesByUser}>Lấy Roles theo User</Button>
        </Row>
        <Row>
          <Field label="Role ID (uuid)">
            <input value={query.role_id} onChange={(e) => setQuery({ ...query, role_id: e.target.value })} placeholder="uuid role" />
          </Field>
          <Button disabled={!notEmpty(query.role_id) || busy} onClick={fetchPermissionsByRole}>Lấy Permissions theo Role</Button>
        </Row>

        {userPermissions?.length ? (
          <>
            <div>Permissions theo User:</div>
            <Code>{JSON.stringify(userPermissions, null, 2)}</Code>
          </>
        ) : null}

        {userRoles?.length ? (
          <>
            <div>Roles theo User:</div>
            <Code>{JSON.stringify(userRoles, null, 2)}</Code>
          </>
        ) : null}

        {rolePermissions?.length ? (
          <>
            <div>Permissions theo Role:</div>
            <Code>{JSON.stringify(rolePermissions, null, 2)}</Code>
          </>
        ) : null}
      </Section>

      <Section title="Người dùng">
        <Row>
          <Field label="Tìm kiếm">
            <input value={userSearch} onChange={(e) => setUserSearch(e.target.value)} placeholder="tên, email, id..." />
          </Field>
          <Button disabled={usersLoading || busy} onClick={onSearchUsers}>Tìm</Button>
          <Field label="Kích thước trang">
            <select value={size} onChange={onChangeSize}>
              <option value={5}>5</option>
              <option value={10}>10</option>
              <option value={20}>20</option>
              <option value={50}>50</option>
            </select>
          </Field>
          <div style={{ marginLeft: 'auto' }}>
            <Button disabled={usersLoading || page <= 1} onClick={() => onChangePage(-1)}>Trang trước</Button>
            <span style={{ margin: '0 8px' }}>Trang {page}</span>
            <Button disabled={usersLoading || (page * size) >= usersTotal} onClick={() => onChangePage(1)}>Trang sau</Button>
          </div>
        </Row>

        {usersLoading ? <div>Đang tải danh sách...</div> : (
          <div className="admin-table-wrap">
            <table className="admin-table">
              <thead>
                <tr>
                  <th>Người dùng</th>
                  <th>Email</th>
                  <th>Trạng thái</th>
                  <th>Vai trò</th>
                  <th>Ngày tạo</th>
                  <th>Hành động</th>
                </tr>
              </thead>
              <tbody>
                {users.map((u) => {
                  const id = u.id || u.user_id || u.uuid || '';
                  const name = u.names || u.name || u.full_name || '';
                  const email = maskEmail(u.email);
                  const role = (u.role_name || (Array.isArray(u.user_type) ? u.user_type.join(', ') : (u.user_type || '')) || '').toString();
                  const createdAt = formatDate(u.created_at);
                  const profilePic = u.profile_pic || '';
                  const stateInfo = getStatePill(u.states);
                  return (
                    <tr key={id}>
                      <td>
                        <div className="admin-row" style={{ alignItems: 'center', gap: 10 }}>
                          {profilePic ? (
                            <img className="admin-avatar" src={profilePic} alt={name || id} />
                          ) : (
                            <span className="admin-avatar-fallback">{(name || id || '?').toString().charAt(0).toUpperCase()}</span>
                          )}
                          <div>
                            <div style={{ fontWeight: 700 }}>{name || '(Chưa có tên)'}</div>
                            <div className="admin-help">{id}</div>
                          </div>
                        </div>
                      </td>
                      <td>{email}</td>
                      <td><span className={stateInfo.cls}>{stateInfo.label}</span></td>
                      <td>{role || '-'}</td>
                      <td>{createdAt}</td>
                      <td>
                        <Button onClick={() => viewUserDetail(id)}>Xem chi tiết</Button>
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        )}

        <div style={{ marginTop: 12 }}>
          {userDetailLoading && <div>Đang tải chi tiết...</div>}
          {userDetailError && <div className="admin-error">{userDetailError}</div>}
          {userDetail && (() => {
            const id = userDetail.id || userDetail.user_id || userDetail.uuid || '';
            const name = userDetail.names || userDetail.name || userDetail.full_name || '';
            const email = userDetail.email || '';
            const role = (userDetail.role_name || (Array.isArray(userDetail.user_type) ? userDetail.user_type.join(', ') : (userDetail.user_type || '')) || '').toString();
            const createdAt = formatDate(userDetail.created_at);
            const profilePic = userDetail.profile_pic || '';
            const stateInfo = getStatePill(userDetail.states);
            return (
              <div className="admin-card" style={{ marginTop: 8 }}>
                <div className="admin-row" style={{ alignItems: 'center', gap: 12 }}>
                  {profilePic ? (
                    <img className="admin-avatar" src={profilePic} alt={name || id} />
                  ) : (
                    <span className="admin-avatar-fallback">{(name || id || '?').toString().charAt(0).toUpperCase()}</span>
                  )}
                  <div>
                    <div style={{ fontWeight: 700, fontSize: 16 }}>{name || '(Chưa có tên)'}</div>
                    <div className="admin-help">{id}</div>
                  </div>
                </div>
                <div className="admin-row" style={{ marginTop: 12 }}>
                  <div className="admin-field" style={{ minWidth: 220 }}>
                    <div className="admin-field-label">Email</div>
                    <div>{email}</div>
                  </div>
                  <div className="admin-field" style={{ minWidth: 160 }}>
                    <div className="admin-field-label">Trạng thái</div>
                    <div><span className={stateInfo.cls}>{stateInfo.label}</span></div>
                  </div>
                  <div className="admin-field" style={{ minWidth: 160 }}>
                    <div className="admin-field-label">Vai trò</div>
                    <div>{role || '-'}</div>
                  </div>
                  <div className="admin-field" style={{ minWidth: 220 }}>
                    <div className="admin-field-label">Ngày tạo</div>
                    <div>{createdAt}</div>
                  </div>
                </div>
              </div>
            );
          })()}
        </div>
      </Section>
    </div>
  );
};

export default AdminDashboardPage;


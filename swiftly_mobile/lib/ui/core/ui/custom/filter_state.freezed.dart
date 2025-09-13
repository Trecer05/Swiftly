// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'filter_state.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

T _$identity<T>(T value) => value;

final _privateConstructorUsedError = UnsupportedError(
  'It seems like you constructed your class using `MyClass._()`. This constructor is only meant to be used by freezed and you are not supposed to need it nor use it.\nPlease check the documentation here for more information: https://github.com/rrousselGit/freezed#adding-getters-and-methods-to-our-models',
);

/// @nodoc
mixin _$FilterState {
  String? get labelTitle => throw _privateConstructorUsedError;
  String? get priorityTitle => throw _privateConstructorUsedError;
  String? get userName => throw _privateConstructorUsedError;

  /// Create a copy of FilterState
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $FilterStateCopyWith<FilterState> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $FilterStateCopyWith<$Res> {
  factory $FilterStateCopyWith(
    FilterState value,
    $Res Function(FilterState) then,
  ) = _$FilterStateCopyWithImpl<$Res, FilterState>;
  @useResult
  $Res call({String? labelTitle, String? priorityTitle, String? userName});
}

/// @nodoc
class _$FilterStateCopyWithImpl<$Res, $Val extends FilterState>
    implements $FilterStateCopyWith<$Res> {
  _$FilterStateCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of FilterState
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? labelTitle = freezed,
    Object? priorityTitle = freezed,
    Object? userName = freezed,
  }) {
    return _then(
      _value.copyWith(
            labelTitle:
                freezed == labelTitle
                    ? _value.labelTitle
                    : labelTitle // ignore: cast_nullable_to_non_nullable
                        as String?,
            priorityTitle:
                freezed == priorityTitle
                    ? _value.priorityTitle
                    : priorityTitle // ignore: cast_nullable_to_non_nullable
                        as String?,
            userName:
                freezed == userName
                    ? _value.userName
                    : userName // ignore: cast_nullable_to_non_nullable
                        as String?,
          )
          as $Val,
    );
  }
}

/// @nodoc
abstract class _$$FilterStateImplCopyWith<$Res>
    implements $FilterStateCopyWith<$Res> {
  factory _$$FilterStateImplCopyWith(
    _$FilterStateImpl value,
    $Res Function(_$FilterStateImpl) then,
  ) = __$$FilterStateImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call({String? labelTitle, String? priorityTitle, String? userName});
}

/// @nodoc
class __$$FilterStateImplCopyWithImpl<$Res>
    extends _$FilterStateCopyWithImpl<$Res, _$FilterStateImpl>
    implements _$$FilterStateImplCopyWith<$Res> {
  __$$FilterStateImplCopyWithImpl(
    _$FilterStateImpl _value,
    $Res Function(_$FilterStateImpl) _then,
  ) : super(_value, _then);

  /// Create a copy of FilterState
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? labelTitle = freezed,
    Object? priorityTitle = freezed,
    Object? userName = freezed,
  }) {
    return _then(
      _$FilterStateImpl(
        labelTitle:
            freezed == labelTitle
                ? _value.labelTitle
                : labelTitle // ignore: cast_nullable_to_non_nullable
                    as String?,
        priorityTitle:
            freezed == priorityTitle
                ? _value.priorityTitle
                : priorityTitle // ignore: cast_nullable_to_non_nullable
                    as String?,
        userName:
            freezed == userName
                ? _value.userName
                : userName // ignore: cast_nullable_to_non_nullable
                    as String?,
      ),
    );
  }
}

/// @nodoc

class _$FilterStateImpl implements _FilterState {
  const _$FilterStateImpl({this.labelTitle, this.priorityTitle, this.userName});

  @override
  final String? labelTitle;
  @override
  final String? priorityTitle;
  @override
  final String? userName;

  @override
  String toString() {
    return 'FilterState(labelTitle: $labelTitle, priorityTitle: $priorityTitle, userName: $userName)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$FilterStateImpl &&
            (identical(other.labelTitle, labelTitle) ||
                other.labelTitle == labelTitle) &&
            (identical(other.priorityTitle, priorityTitle) ||
                other.priorityTitle == priorityTitle) &&
            (identical(other.userName, userName) ||
                other.userName == userName));
  }

  @override
  int get hashCode =>
      Object.hash(runtimeType, labelTitle, priorityTitle, userName);

  /// Create a copy of FilterState
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$FilterStateImplCopyWith<_$FilterStateImpl> get copyWith =>
      __$$FilterStateImplCopyWithImpl<_$FilterStateImpl>(this, _$identity);
}

abstract class _FilterState implements FilterState {
  const factory _FilterState({
    final String? labelTitle,
    final String? priorityTitle,
    final String? userName,
  }) = _$FilterStateImpl;

  @override
  String? get labelTitle;
  @override
  String? get priorityTitle;
  @override
  String? get userName;

  /// Create a copy of FilterState
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$FilterStateImplCopyWith<_$FilterStateImpl> get copyWith =>
      throw _privateConstructorUsedError;
}
